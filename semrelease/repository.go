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
	initialVersion = "v1.0.0"
)

// Repository ...
type Repository interface {
	CloneRepository(ctx context.Context, owner, repo, token string) error
	ListCommits(ctx context.Context, owner, repo string, currentVersion string) ([]Commit, error)
	GetLatestVersion(ctx context.Context, owner, repo string) (string, error)

	getReference(ctx context.Context, owner, repo, reference string) (*github.Reference, error)
	createRelease(ctx context.Context, owner, repo string, repositoryRelease *github.RepositoryRelease) (*github.RepositoryRelease, error)
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
func (r GitHubRepository) CloneRepository(ctx context.Context, owner, repo, token string) error {
	url := fmt.Sprintf("https://%s:x-oauth-basic@github.com/%s/%s.git", token, owner, repo)
	cmd := exec.CommandContext(ctx, "git", "clone", url)
	return cmd.Run()
}

// ListCommits ...
func (r GitHubRepository) ListCommits(ctx context.Context, owner, repo string, currentVersion string) ([]Commit, error) {
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
	jsonCommits := fmt.Sprintf("[%s]", output[:len(output)-2])
	err = json.Unmarshal([]byte(jsonCommits), &commits)

	return commits, err
}

// GetLatestVersion ...
func (r GitHubRepository) GetLatestVersion(ctx context.Context, owner, repo string) (string, error) {
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
	// if version == "" {
	// 	version = initialVersion
	// }
	return version, nil
}

func (r GitHubRepository) getReference(ctx context.Context, owner, repo, reference string) (*github.Reference, error) {
	ref, _, err := r.Client.Git.GetRef(ctx, owner, repo, reference)
	return ref, err
}

func (r GitHubRepository) createRelease(ctx context.Context, owner, repo string, repositoryRelease *github.RepositoryRelease) (*github.RepositoryRelease, error) {
	release, _, err := r.Client.Repositories.CreateRelease(ctx, owner, repo, repositoryRelease)
	return release, err
}
