package main

import (
	"context"
	"fmt"

	"github.com/catho/Sem-Release/semrelease"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	config := newConfig()

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.AccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	repository := semrelease.NewRepository(client)
	service := semrelease.NewService(repository)
	version, err := service.CreateRelease(ctx, config.Owner, config.Repository, config.AccessToken, config.ReleaseBranch)
	if err != nil {
		panic(err)
	}

	fmt.Printf(fmt.Sprintf(`{"version": "v%v"}`, version))

}
