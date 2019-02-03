package semrelease

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"

	"github.com/Masterminds/semver"
)

type Service struct {
	GitRepositoryApplication IRepository
}

// CalculateChanges ...
func (s Service) CalculateChanges(commits []*Commit, latestRelease *Release) Change {
	var change Change
	for _, commit := range commits {
		if latestRelease.SHA == commit.SHA {
			break
		}
		change.Major = change.Major || commit.Change.Major
		change.Minor = change.Minor || commit.Change.Minor
		change.Patch = change.Patch || commit.Change.Patch
	}
	return change
}

// GetNewVersion calcula new version
func (s Service) GetNewVersion(commits []*Commit, latestRelease *Release) *semver.Version {
	if latestRelease == nil {
		return s.ApplyChange(&semver.Version{}, Change{})
	}
	return s.ApplyChange(latestRelease.Version, s.CalculateChanges(commits, latestRelease))
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
	"feat":     "Feature",
	"fix":      "Bug Fixes",
	"perf":     "Performance Improvements",
	"revert":   "Reverts",
	"docs":     "Documentation",
	"style":    "Styles",
	"refactor": "Code Refactoring",
	"test":     "Tests",
	"chore":    "Chores",
	"%%bc%%":   "Breaking Changes",
}

func formatCommit(c *Commit) string {
	ret := "* "
	if c.Scope != "" {
		ret += fmt.Sprintf("**%s:** ", c.Scope)
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
func (s Service) CreateRelease(owner, repo string) {
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
	lastedRelease, _ := s.GetLatestRelease()
	version := s.GetNewVersion(commits, lastedRelease)
	changelog := s.GetChangelog(commits, lastedRelease, version)
	s.createRelease(owner, repo, changelog, version, false, "master")
}

func (s Service) createRelease(owner, repo, changelog string, newVersion *semver.Version, prerelease bool, branch string) error {

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

	_, _, err := cl.Repositories.CreateRelease(ctx, owner, repo, opts)
	fmt.Println(err)
	if err != nil {
		return err
	}
	return nil
}

// GetLatestRelease ..
func (s Service) GetLatestRelease() (*Release, error) {
	ctx := context.TODO()
	owner := "eduardokenjimiura"
	repo := "REPOSITORIOAPP"
	allReleases := make(Releases, 0)
	opts := &github.ReferenceListOptions{"tags", github.ListOptions{PerPage: 100}}
	cl := s.ReturnClient()
	for {
		refs, resp, err := cl.Git.ListRefs(ctx, owner, repo, opts)
		if resp != nil && resp.StatusCode == 404 {
			return &Release{"", &semver.Version{}}, nil
		}
		if err != nil {
			return nil, err
		}
		for _, r := range refs {
			version, err := semver.NewVersion(strings.TrimPrefix(r.GetRef(), "refs/tags/"))
			if err != nil {
				continue
			}
			allReleases = append(allReleases, &Release{r.Object.GetSHA(), version})
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	var lastRelease *Release
	for _, r := range allReleases {
		if r.Version.Prerelease() == "" {
			lastRelease = r
			break
		}
	}

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
	tr := http.DefaultTransport
	itr, err := ghinstallation.NewKeyFromFile(tr, 22198, 515453, "privatekey.pem")
	if err != nil {
		log.Fatal(err)
	}
	client := github.NewClient(&http.Client{Transport: itr})
	return client
}

var breakingPattern = regexp.MustCompile("BREAKING CHANGES?")

var commitPattern = regexp.MustCompile("^(\\w*)(?:\\((.*)\\))?\\: (.*)$")

func parseCommit(commit *github.RepositoryCommit) *Commit {
	c := new(Commit)
	c.SHA = commit.GetSHA()
	fmt.Println(c.SHA)
	c.Raw = strings.Split(commit.Commit.GetMessage(), "\n")
	strings.Split(commit.Commit.GetMessage(), "\n")

	// for _, sp := range split {
	// 	found := commitPattern.FindAllStringSubmatch(sp, -1)
	// 	if len(found) < 1 {
	// 		fmt.Println(found, "  ==================== ", sp)
	// 		return c
	// 	}
	// 	break
	// }

	found := commitPattern.FindAllStringSubmatch(c.Raw[0], -1)
	if len(found) < 1 {
		return c
	}
	//c.Type = strings.ToLower(found[0][0])

	message := c.Raw[0]
	tp := message[:strings.IndexByte(message, ':')]
	c.Type = tp
	if len(tp) > 0 {
		scope := message[strings.IndexByte(message, ':'):]
		c.Scope = scope
		//	c.Message = found[0][3]
	}
	c.Change = Change{
		Major: breakingPattern.MatchString(c.Raw[0]),
		Minor: c.Type == "feat",
		Patch: c.Type == "fix",
	}
	return c
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
