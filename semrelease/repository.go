package semrelease

import (
	"context"

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
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "15bcf142682db5244445b83da33b1f872a399624"},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return &RepositoryImpl{Client: *client}
}

func (r RepositoryImpl) AllRepository(ctx context.Context) ([]*github.Repository, error) {
	repos, _, err := r.Client.Repositories.List(ctx, "", nil)
	return repos, err
}
