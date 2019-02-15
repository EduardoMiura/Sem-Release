package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	AccessToken   string `envconfig:"ACCESS_TOKEN" required:"true"`
	Owner         string `envconfig:"OWNER" required:"true"`
	Repository    string `envconfig:"REPOSITORY" required:"true"`
	ReleaseBranch string `envconfig:"RELEASE_BRANCH" default:"master"`
}

func newConfig() config {
	cfg := &config{}
	if err := envconfig.Process("", cfg); err != nil {
		log.Fatal(err)
	}
	return *cfg
}
