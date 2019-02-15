package semrelease

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/google/go-github/github"
)

const (
	logFormat      = `{"sha": "%H", "abbreviatedSHA": "%h", "message": "%s", "sanitizedMessage": "%f"},`
	initialVersion = "vT1.0.0"
)

// Repository ...
type Repository interface {
	cloneRepository(ctx context.Context, owner, repo, token string) error
	listCommits(ctx context.Context, owner, repo string, currentVersion string) ([]Commit, error)
	getLatestVersion(ctx context.Context, owner, repo string) (string, error)
	createRelease(ctx context.Context, release Release) (*github.RepositoryRelease, error)
}

// GitHubRepository ...
type GitHubRepository struct {
	Client *github.Client
}

// NewRepository ...
func NewRepository(client *github.Client) *GitHubRepository {
	return &GitHubRepository{
		Client: client,
	}
}

// CloneRepository ...
func (r GitHubRepository) cloneRepository(ctx context.Context, owner, repo, token string) error {
	url := fmt.Sprintf("https://%s:x-oauth-basic@github.com/%s/%s.git", token, owner, repo)
	cmd := exec.CommandContext(ctx, "git", "clone", url)
	return cmd.Run()
}

// ListCommits ...
func (r GitHubRepository) listCommits(ctx context.Context, owner, repo string, currentVersion string) ([]Commit, error) {
	var commits []Commit
	var stdout bytes.Buffer

	path, err := os.Getwd()
	if err != nil {
		return commits, err
	}

	args := []string{
		"log",
		fmt.Sprintf("--format=%s", logFormat),
	}
	if currentVersion != "" {
		args = append(args, fmt.Sprintf("%s..", currentVersion))
	}

	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = fmt.Sprintf("%s/%s", path, repo)
	cmd.Stdout = &stdout
	err = cmd.Run()
	if err != nil {
		return commits, err
	}

	output := stdout.String()
	
	if output != "" {
		jsonCommits := fmt.Sprintf("[%s]", output[:len(output)-2])
		err = json.Unmarshal([]byte(jsonCommits), &commits)
	}
	return commits, err
}

// GetLatestVersion ...
func (r GitHubRepository) getLatestVersion(ctx context.Context, owner, repo string) (string, error) {
	var stdout bytes.Buffer
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}

	cmd := exec.CommandContext(ctx, "git", "for-each-ref", "refs/tags", "--sort=-taggerdate", "--format=%(refname:short)", "--count=1")
	cmd.Dir = fmt.Sprintf("%s/%s", path, repo)
	cmd.Stdout = &stdout
	err = cmd.Run()
	if err != nil {
		return "", err
	}

	version := strings.Trim(stdout.String(), "\n")
	return version, nil
}

func (r GitHubRepository) createRelease(ctx context.Context, release Release) (*github.RepositoryRelease, error) {
	tag := fmt.Sprintf("vT%s", release.Version.String())
	releaseNote := release.getReleaseNote()
	repositoryRelease := &github.RepositoryRelease{
		TagName:         &tag,
		Name:            &tag,
		TargetCommitish: &release.Branch,
		Body:            &releaseNote,
		Prerelease:      &release.IsPreRelease,
	}

	createdRelease, _, err := r.Client.Repositories.CreateRelease(ctx, release.Owner, release.Repository, repositoryRelease)
	return createdRelease, err
}
