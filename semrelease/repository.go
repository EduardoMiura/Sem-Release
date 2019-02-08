package semrelease

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/google/go-github/github"
)

const (
	logFormat = `{"sha": "%H", "abbreviatedSHA": "%h", "message": "%s", "sanitizedMessage": "%f"},`
)

// Repository ...
type Repository interface {
	CloneRepository(ctx context.Context, owner, repo, token string) error

	listCommits(ctx context.Context, owner, repo string, latestRelease *Release) ([]*github.RepositoryCommit, error)
	getLatestRelease(ctx context.Context, owner, repo string) (*github.RepositoryRelease, error)
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

// TODO: listar a partir do ultimo release
func (r GitHubRepository) listCommits(ctx context.Context, owner, repo string, lastestRelease *Release) ([]Commit, error) {
	var commits []Commit
	var stdout bytes.Buffer

	version := lastestRelease.Version.String()
	path, err := os.Getwd()
	if err != nil {
		return commits, err
	}

	cmd := exec.CommandContext(ctx, "git", "log", fmt.Sprintf("v%s..", version), fmt.Sprintf("--format=%s", logFormat))
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

func (r GitHubRepository) getLatestRelease(ctx context.Context, owner, repo string) (*github.RepositoryRelease, error) {


//git for-each-ref refs/tags --sort=-taggerdate --format='%(refname)' --count=1

	latestRelease, _, err := r.Client.Repositories.GetLatestRelease(ctx, owner, repo)
	return latestRelease, err
}

func (r GitHubRepository) getReference(ctx context.Context, owner, repo, reference string) (*github.Reference, error) {
	ref, _, err := r.Client.Git.GetRef(ctx, owner, repo, reference)
	return ref, err
}

func (r GitHubRepository) createRelease(ctx context.Context, owner, repo string, repositoryRelease *github.RepositoryRelease) (*github.RepositoryRelease, error) {
	release, _, err := r.Client.Repositories.CreateRelease(ctx, owner, repo, repositoryRelease)
	return release, err
}
