package main

import (
	"context"
	"os"

	"github.com/catho/Sem-Release/semrelease"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	accessToken := os.Getenv("ACCESS_TOKEN")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	repository := semrelease.NewRepository(client)
	service := semrelease.NewService(repository)
	owner := os.Getenv("OWNER")
	repo := os.Getenv("REPOSITORY")
	service.CreateRelease(ctx, owner, repo)
	repository.CloneRepository(ctx, owner, repo, accessToken)
}

//url := fmt.Sprintf("https://%s:x-oauth-basic@github.com/%s/%s.git", token, owner, repo)
