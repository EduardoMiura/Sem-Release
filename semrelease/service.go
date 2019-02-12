package semrelease

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"strings"
	"time"

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
func (s Service) CreateRelease(ctx context.Context, owner, repo string) (*semver.Version, error) {
	latestRelease, err := s.Repository.GetLatestRelease(ctx, owner, repo) //s.getLatestRelease(ctx, owner, repo)
	if err != nil {
		return nil, err
	}

	cm, err := s.Repository.ListCommits(ctx, owner, repo, latestRelease)
	if err != nil {
		return nil, err
	}

	// var commits []*Commit
	// for i, commit := range cm {
	// 	c := parseCommit(commit)
	// 	if c != nil {
	// 		commits = append(commits, c)
	// 	}

	// 	b, _ := json.Marshal(c)
	// 	log.Println("parse", i, string(b))
	// }

	s.newRelease(ctx, owner, repo, "master", latestRelease, cm)

	// version := s.getNewVersion(commits, lastedRelease)
	// changelog := s.getChangelog(commits, lastedRelease, version)
	// rep, err := s.createRelease(ctx, owner, repo, changelog, version, false, "master")
	// createFileRelease(rep)
	// return version, err
	return nil, err
}

func (s Service) newRelease(ctx context.Context, owner, repo, branch, currentRelease string, commits []Commit) Release {
	release := Release{
		ChangeLog: map[string][]string{},
		SHA:       commits[0].SHA,
		Branch:    branch,
	}

	for _, commit := range commits {
		c := parseCommit(commit)
		if c == nil {
			continue
		}
		commit = *c
		release.ChangeLog[commit.Type.Name] = append(release.ChangeLog[commit.Type.Name], fmt.Sprintf("* %s: %s (%s)", commit.Scope, commit.Message, commit.AbbreviatedSHA))

		commitType := commitTypes[commit.Type.Name]
		if commitType.priority.value > release.Change.value {
			release.Change = commitType.priority
		}
	}

	version, _ := semver.NewVersion(currentRelease)
	var v semver.Version
	switch release.Change.Name {
	case "major":
		v = version.IncMajor()
	case "minor":
		v = version.IncMinor()
	case "patch":
		v = version.IncPatch()
	}

	release.Version = &v
	fmt.Println("AAA", currentRelease, version.String(), release.Version.String())
	fmt.Println("CHANGE-LOG", release.ChangeLog.String())

	b, _ := json.Marshal(release)
	fmt.Println("RELEASE", string(b))

	s.createRelease(ctx, owner, repo, release.ChangeLog.String(), release.Version, false, "master")

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

// isPatch ...
func isPatch(typeOfChange string) bool {
	switch typeOfChange {
	case
		"fix",
		"perf",
		"revert",
		"docs",
		"style",
		"refactor",
		"test",
		"chore":
		return true
	}
	return false
}

// calculateChanges ...
func (s Service) calculateChanges(commits []*Commit, latestRelease *Release) Change {
	var change Change
	for _, commit := range commits {
		change.Major = commit.Change.Major
		if change.Major {
			break
		}

		change.Minor = commit.Change.Minor
		change.Patch = commit.Change.Patch
	}
	return change
}

// getNewVersion create new version
// func (s Service) getNewVersion(commits []*Commit, latestRelease *Release) *semver.Version {
// 	if latestRelease == nil {
// 		return s.applyChange(&semver.Version{}, Change{})
// 	}
// 	ch := s.calculateChanges(commits, latestRelease)
// 	return s.applyChange(latestRelease.Version, ch)
// }

// // applyChange ...
// func (s Service) applyChange(version *semver.Version, change Change) *semver.Version {
// 	if version.Major() == 0 {
// 		change.Major = true
// 	}
// 	if !change.Major && !change.Minor && !change.Patch {
// 		return version
// 	}
// 	var newVersion semver.Version
// 	preRel := version.Prerelease()

// 	if preRel == "" {
// 		switch {
// 		case change.Major:
// 			newVersion = version.IncMajor()
// 			break
// 		case change.Minor:
// 			newVersion = version.IncMinor()
// 			break
// 		case change.Patch:
// 			newVersion = version.IncPatch()
// 			break
// 		}
// 		return &newVersion
// 	}
// 	preRelVer := strings.Split(preRel, ".")
// 	if len(preRelVer) > 1 {
// 		idx, err := strconv.ParseInt(preRelVer[1], 10, 32)
// 		if err != nil {
// 			idx = 0
// 		}
// 		preRel = fmt.Sprintf("%s.%d", preRelVer[0], idx+1)
// 	} else {
// 		preRel += ".1"
// 	}
// 	newVersion, _ = version.SetPrerelease(preRel)
// 	return &newVersion
// }

var typeToText = map[string]string{
	"feat":     "A new Feature",
	"fix":      "A Bug Fixes",
	"perf":     "Performance Improvements",
	"revert":   "Reverts",
	"docs":     "Documentation only change",
	"style":    "Styles",
	"refactor": "Code Refactoring",
	"test":     "Add Tests",
	"chore":    "Change to the build process Chores",
	"%%bc%%":   "Breaking Changes",
}

func trimSHA(sha string) string {
	if len(sha) < 9 {
		return sha
	}
	return sha[:8]
}
func getSortedKeys(m *map[string]string) []string {
	keys := make([]string, len(*m))
	i := 0
	for k := range *m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
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

// getChangelog ..
func (s Service) getChangelog(commits []*Commit, latestRelease *Release, newVersion *semver.Version) string {
	ret := fmt.Sprintf("## %s (%s)\n\n", newVersion.String(), time.Now().UTC().Format("2006-01-02"))
	typeScopeMap := make(map[string]string)
	for _, commit := range commits {
		if latestRelease.SHA == commit.SHA {
			break
		}
		if commit.Change.Major {
			typeScopeMap["%%bc%%"] += fmt.Sprintf("%s\n```%s\n```\n", formatCommit(commit), strings.Join(commit.Raw[1:], "\n"))
			continue
		}
		if commit.Type.Name == "" {
			continue
		}
		typeScopeMap[commit.Type.Name] += formatCommit(commit)
	}
	for _, t := range getSortedKeys(&typeScopeMap) {
		msg := typeScopeMap[t]
		typeName, found := typeToText[t]
		if !found {
			typeName = t
		}
		ret += fmt.Sprintf("#### %s\n\n%s\n", typeName, msg)
	}
	return ret
}

func formatCommit(c *Commit) string {
	ret := "* "
	if c.Scope != "" {
		ret += fmt.Sprintf("%s: ", c.Scope)
	}
	ret += fmt.Sprintf("%s (%s)\n", c.Message, trimSHA(c.SHA))
	return ret
}
