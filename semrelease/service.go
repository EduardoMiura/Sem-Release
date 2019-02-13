package semrelease

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/google/go-github/github"

	"github.com/Masterminds/semver"
)

var (
	breakingPattern = regexp.MustCompile("BREAKING CHANGES?")
	commitPattern   = regexp.MustCompile(`^(?P<type>[\w\s]*)(?:\((?P<scope>.*)\))?\s*\:(?P<message>.*)$`)
)

// Service ...
type Service struct {
	Repository Repository
}

// NewService ...
func NewService(repository Repository) Service {
	return Service{
		Repository: repository,
	}
}

// CreateRelease ...
func (s Service) CreateRelease(ctx context.Context, owner, repo, releaseBranch string) (*semver.Version, error) {
	currentVersion, err := s.Repository.GetLatestVersion(ctx, owner, repo)
	if err != nil {
		return nil, err
	}

	cm, err := s.Repository.ListCommits(ctx, owner, repo, currentVersion)
	if err != nil {
		return nil, err
	}

	release := s.newRelease(ctx, owner, repo, releaseBranch, currentVersion, cm)
	s.createRelease(ctx, owner, repo, release.getReleaseNote(), release.Version, false, "master")

	return nil, err
}

func (s Service) newRelease(ctx context.Context, owner, repo, branch, currentVersion string, commits []Commit) Release {
	var previousVersion *semver.Version
	if currentVersion != "" {
		previousVersion, _ = semver.NewVersion(currentVersion)
	}

	release := Release{
		ChangeLog:       map[string][]string{},
		SHA:             commits[0].SHA,
		Branch:          branch,
		PreviousVersion: previousVersion,
		Owner:           owner,
		Repository:      repo,
	}

	for _, commit := range commits {
		c := parseCommit(commit)
		if c == nil {
			continue
		}
		commit = *c
		release.ChangeLog[commit.Type.Name] = append(release.ChangeLog[commit.Type.Name], fmt.Sprintf("* **%s:** %s (%s)", commit.Scope, commit.Message, commit.AbbreviatedSHA))

		commitType := commitTypes[commit.Type.Name]
		if commitType.priority.value > release.Change.value {
			release.Change = commitType.priority
		}
	}

	var newVersion semver.Version

	if previousVersion != nil {
		switch release.Change.Name {
		case "major":
			newVersion = release.PreviousVersion.IncMajor()
		case "minor":
			newVersion = release.PreviousVersion.IncMinor()
		case "patch":
			newVersion = release.PreviousVersion.IncPatch()
		}
		release.Version = &newVersion
	} else {
		release.Version, _ = semver.NewVersion("1.0.0")
	}

	if release.PreviousVersion != nil {

		fmt.Println("AAA", currentVersion, release.PreviousVersion.String(), release.Version.String())
	} else {
		fmt.Println("ABA", currentVersion, release.Version.String())
	}

	fmt.Println("CHANGE-LOG", release.ChangeLog.String())

	b, _ := json.Marshal(release)
	fmt.Println("RELEASE", string(b))

	return release
}

var commitTypes = map[string]CommitType{
	"breaking": CommitType{Name: "breaking", Description: "Breaking Changes", priority: major},
	"bc":       CommitType{Name: "bc", Description: "Breaking Changes", priority: major},
	"feat":     CommitType{Name: "feat", Description: "A new Feature", priority: minor},
	"fix":      CommitType{Name: "fix", Description: "A Bug Fixes", priority: patch},
	"perf":     CommitType{Name: "perf", Description: "Performance Improvements", priority: patch},
	"revert":   CommitType{Name: "revert", Description: "Reverts", priority: patch},
	"docs":     CommitType{Name: "docs", Description: "Documentation only change", priority: patch},
	"style":    CommitType{Name: "style", Description: "Styles", priority: patch},
	"refactor": CommitType{Name: "refactor", Description: "Code Refactoring", priority: patch},
	"test":     CommitType{Name: "test", Description: "Add Tests", priority: patch},
	"chore":    CommitType{Name: "chore", Description: "Change to the build process Chores", priority: patch},
}

//func parseCommit(commit *github.RepositoryCommit) *Commit {
func parseCommit(commit Commit) *Commit {
	message := strings.TrimSpace(commit.Message)
	if !commitPattern.MatchString(message) {
		return nil
	}

	commitValues := map[string]string{}
	groupNames := commitPattern.SubexpNames()
	matches := commitPattern.FindStringSubmatch(message)

	for i, name := range groupNames {
		if name == "" {
			continue
		}
		commitValues[name] = strings.Trim(matches[i], " ")
	}

	// c := new(Commit)
	commit.Raw = strings.Split(message, " ")
	commitType := commitValues["type"]
	commit.Type = commitTypes[commitType]
	commit.Scope = commitValues["scope"]
	commit.Message = commitValues["message"]

	return &commit
}

func createFileRelease(data interface{}) {
	b, _ := json.MarshalIndent(data, "", " ")
	fmt.Println(string(b))
	er := ioutil.WriteFile("file.json", b, 0644)
	if er != nil {
		fmt.Println("file is generate ")
	}
}

func (s Service) createRelease(ctx context.Context, owner, repo, changelog string, newVersion *semver.Version, prerelease bool, branch string) (*github.RepositoryRelease, error) {
	tag := fmt.Sprintf("v%s", newVersion.String())
	isPrerelease := prerelease || newVersion.Prerelease() != ""
	opts := &github.RepositoryRelease{
		TagName:         &tag,
		Name:            &tag,
		TargetCommitish: &branch,
		Body:            &changelog,
		Prerelease:      &isPrerelease,
	}

	return s.Repository.createRelease(ctx, owner, repo, opts)
}
