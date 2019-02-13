package main

import (
	"context"
	"log"
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
	owner := os.Getenv("OWNER")
	repo := os.Getenv("REPOSITORY")

	err := repository.CloneRepository(ctx, owner, repo, accessToken)
	if err != nil {
		log.Fatal("clone_repository", err)
	}

	releaseBranch := os.Getenv("RELEASE_BRANCH")
	service := semrelease.NewService(repository)
	service.CreateRelease(ctx, owner, repo, releaseBranch)

	// TODO: update repository to make funcs private. (they are public to test repository)
	version, err := repository.GetLatestVersion(ctx, owner, repo)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("version...", version)
	commits, err := repository.ListCommits(ctx, owner, repo, version)
	if err != nil {
		log.Fatal("commits ", err, commits)
	}
}

/*TODO:

- parametrizar qual é a release branch
- utilizar template para gerar o change log
-

*/
