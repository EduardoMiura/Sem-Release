package semrelease

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"github.com/Masterminds/semver"
)

type Service struct {
	GitRepositoryApplication IRepository
}

// CalculateChanges ...
func (s Service) CalculateChanges(commits []*Commit, latestRelease *Release) Change {
	var change Change

	for _, commit := range commits {
		change.Major = commit.Change.Major
		change.Minor = commit.Change.Minor
		change.Patch = commit.Change.Patch
		break

	}

	return change
}

// GetNewVersion calcula new version
func (s Service) GetNewVersion(commits []*Commit, latestRelease *Release) *semver.Version {
	if latestRelease == nil {
		return s.ApplyChange(&semver.Version{}, Change{})
	}
	fmt.Println(commits[0], "commits")
	ch := s.CalculateChanges(commits, latestRelease)
	return s.ApplyChange(latestRelease.Version, ch)
}

// CheckHealth ...
func (s Service) CheckHealth(ctx context.Context, timeout time.Duration) error {

	return nil
}

// ApplyChange ...
func (s Service) ApplyChange(version *semver.Version, change Change) *semver.Version {
	if version.Major() == 0 {
		change.Major = true
	}
	if !change.Major && !change.Minor && !change.Patch {
		return version
	}
	var newVersion semver.Version
	preRel := version.Prerelease()

	if preRel == "" {
		switch {
		case change.Major:
			newVersion = version.IncMajor()
			break
		case change.Minor:
			newVersion = version.IncMinor()
			break
		case change.Patch:
			newVersion = version.IncPatch()
			break
		}
		return &newVersion
	}
	preRelVer := strings.Split(preRel, ".")
	if len(preRelVer) > 1 {
		idx, err := strconv.ParseInt(preRelVer[1], 10, 32)
		if err != nil {
			idx = 0
		}
		preRel = fmt.Sprintf("%s.%d", preRelVer[0], idx+1)
	} else {
		preRel += ".1"
	}
	newVersion, _ = version.SetPrerelease(preRel)
	return &newVersion
}

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

func formatCommit(c *Commit) string {
	ret := "* "
	if c.Scope != "" {
		ret += fmt.Sprintf("%s: ", c.Scope)
	}
	ret += fmt.Sprintf("%s (%s)\n", c.Message, trimSHA(c.SHA))
	return ret
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

// CreateRelease ...
func (s Service) CreateRelease(owner, repo string) *semver.Version {

	cl := s.ReturnClient()

	ctx := context.TODO()
	cm, _, _ := cl.Repositories.ListCommits(ctx, owner, repo, nil)
	var commits []*Commit
	for _, commit := range cm {
		c := parseCommit(commit)

		if c.Type != "" {
			commits = append(commits, c)
			break
		}
	}
	lastedRelease, _ := s.GetLatestRelease(owner, repo)

	version := s.GetNewVersion(commits, lastedRelease)

	changelog := s.GetChangelog(commits, lastedRelease, version)
	_, rep := s.createRelease(owner, repo, changelog, version, false, "master")
	createFileRelease(rep)
	return version
}

func createFileRelease(data interface{}) {
	b, _ := json.MarshalIndent(data, "", " ")
	fmt.Println(string(b))
	er := ioutil.WriteFile("file.json", b, 0644)
	if er != nil {
		fmt.Println("file is generate ")
	}
}

func (s Service) createRelease(owner, repo, changelog string, newVersion *semver.Version, prerelease bool, branch string) (error, *github.RepositoryRelease) {
	tag := fmt.Sprintf("v%s", newVersion.String())
	isPrerelease := prerelease || newVersion.Prerelease() != ""
	ctx := context.TODO()

	opts := &github.RepositoryRelease{
		TagName:         &tag,
		Name:            &tag,
		TargetCommitish: &branch,
		Body:            &changelog,
		Prerelease:      &isPrerelease,
	}
	cl := s.ReturnClient()
	repreturn, _, err := cl.Repositories.CreateRelease(ctx, owner, repo, opts)
	fmt.Println(err)
	if err != nil {
		return err, nil
	}
	return err, repreturn
}

// GetLatestRelease ..
func (s Service) GetLatestRelease(owner, repo string) (*Release, error) {
	ctx := context.TODO()
	cl := s.ReturnClient()
	RepositoryRelease, _, _ := cl.Repositories.GetLatestRelease(ctx, owner, repo)
	version, _ := semver.NewVersion(RepositoryRelease.GetTagName())
	tagref := "tags/" + RepositoryRelease.GetTagName()
	d, _, _ := cl.Git.GetRef(ctx, owner, repo, tagref)
	lastRelease := &Release{}
	lastRelease.SHA = d.Object.GetSHA()
	lastRelease.Version = version
	return lastRelease, nil
}

// GetChangelog ..
func (s Service) GetChangelog(commits []*Commit, latestRelease *Release, newVersion *semver.Version) string {
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
		if commit.Type == "" {
			continue
		}
		typeScopeMap[commit.Type] += formatCommit(commit)
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
func (s Service) CaptureRepositoryAndOwner(pullRequest interface{}) (string, string) {
	///mudar esse cara para recursivo
	if (pullRequest) != nil {
		head := pullRequest.(map[string]interface{})
		if head != nil {

			repo := head["head"]
			if repo != nil {
				rp := repo.(map[string]interface{})
				if rp != nil {
					owner := rp["repo"]
					if owner != nil {

						ow := owner.(map[string]interface{})
						repositoryname := ow["name"].(string)
						login := ow["owner"].(map[string]interface{})["login"].(string)
						return repositoryname, login
					}
				}
			}

		}
	}
	return "", ""
}

func (s Service) ReturnClient() *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("AccessToken")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client
}

var breakingPattern = regexp.MustCompile("BREAKING CHANGES?")

var commitPattern = regexp.MustCompile("^(\\w*) (?:\\((.*)\\))?\\: (.*)$")

func parseCommit(commit *github.RepositoryCommit) *Commit {
	message := strings.TrimSpace(commit.Commit.GetMessage())
	c := new(Commit)
	c.SHA = commit.GetSHA()
	c.Raw = strings.Split(message, " ")

	// for _, sp := range split {
	// 	found := commitPattern.FindAllStringSubmatch(sp, -1)
	// 	if len(found) < 1 {
	// 		fmt.Println(found, "  ==================== ", sp)
	// 		return c
	// 	}
	// 	break
	// }
	found := commitPattern.MatchString(message)

	if found {
		return c
	}
	tp := c.Raw[0]
	c.Type = tp
	if len(tp) > 0 {
		scope := message[strings.IndexByte(message, ':')+1:]
		c.Scope = scope
		c.Message = c.Raw[1]
	}
	c.Change = Change{
		Major: breakingPattern.MatchString(c.Raw[0]),
		Minor: c.Type == "feat",
		Patch: Patch(c.Type),
	}
	return c
}

func Patch(typeOfChange string) bool {
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

func (s Service) GetRepositories() ([]*github.Repository, error) {

	ctx := context.TODO()
	return s.GitRepositoryApplication.AllRepository(ctx)
}

func NewService(gitRepository IRepository) *Service {
	return &Service{
		GitRepositoryApplication: gitRepository,
	}
}
