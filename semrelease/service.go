package semrelease

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/Masterminds/semver"
)

var (
	commitPattern = regexp.MustCompile(`^(?P<type>[\w\s]*)(?:\((?P<scope>.*)\))?\s*\:(?P<message>.*)$`)
	commitTypes   = map[string]CommitType{
		"breaking": CommitType{Name: "breaking", Description: "Breaking Changes", priority: major},
		"break":    CommitType{Name: "breaking", Description: "Breaking Changes", priority: major},
		"bc":       CommitType{Name: "bc", Description: "Breaking Changes", priority: major},
		"feat":     CommitType{Name: "feat", Description: "Features", priority: minor},
		"fix":      CommitType{Name: "fix", Description: "Bug Fixes", priority: patch},
		"perf":     CommitType{Name: "perf", Description: "Performance Improvements", priority: patch},
		"revert":   CommitType{Name: "revert", Description: "Reverts", priority: patch},
		"docs":     CommitType{Name: "docs", Description: "Documentation", priority: patch},
		"style":    CommitType{Name: "style", Description: "Styles", priority: patch},
		"refactor": CommitType{Name: "refactor", Description: "Code Refactoring", priority: patch},
		"test":     CommitType{Name: "test", Description: "Tests", priority: patch},
		"chore":    CommitType{Name: "chore", Description: "Change to the build process", priority: patch},
	}
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
func (s Service) CreateRelease(ctx context.Context, owner, repo, accessToken, releaseBranch string) (*semver.Version, error) {
	err := s.Repository.cloneRepository(ctx, owner, repo, accessToken)
	if err != nil {
		return nil, err
	}
	err = s.Repository.checkoutBranch(ctx, repo, releaseBranch)
	if err != nil {
		return nil, err
	}
	currentVersion, err := s.Repository.getLatestVersion(ctx, repo)
	if err != nil {
		return nil, err
	}

	cm, err := s.Repository.listCommits(ctx, repo, currentVersion)
	if err != nil {

		return nil, err
	}

	if len(cm) == 0 {
		log.Println("no commits")
		return nil, nil
	}

	release := s.newRelease(ctx, owner, repo, releaseBranch, currentVersion, cm)
	s.Repository.createRelease(ctx, release)

	return release.Version, err
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
		release.ChangeLog[commit.Type.Name] = append(release.ChangeLog[commit.Type.Name], commit.String())

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

	return release
}

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

	commit.Raw = strings.Split(message, " ")
	commitType := commitValues["type"]
	var found bool
	commit.Type, found = commitTypes[commitType]
	if !found {
		return nil
	}
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
