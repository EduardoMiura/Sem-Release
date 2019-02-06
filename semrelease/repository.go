package semrelease

import (
	"context"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type IRepository interface {
	AllRepository(ctx context.Context) ([]*github.Repository, error)
}
type RepositoryImpl struct {
	Client github.Client
}

func NewRepository() *RepositoryImpl {
	ctx := context.Background()
	AccessToken := os.Getenv("AccessToken")
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: AccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return &RepositoryImpl{Client: *client}
}

func (r RepositoryImpl) AllRepository(ctx context.Context) ([]*github.Repository, error) {
	repos, _, err := r.Client.Repositories.List(ctx, "", nil)
	return repos, err
}
