package main

import (
	"flavioltonon/go-chronos/chronos"
	"log"
	"os"

	"github.com/flavioltonon/go-github/github"
)

func main() {
	var c chronos.Chronos

	auth := github.BasicAuthTransport{
		Username: os.Getenv("CHRONOS_GITHUB_LOGIN"),
		Password: os.Getenv("CHRONOS_GITHUB_PASSWORD"),
	}

	c.SetClient(github.NewClient(auth.Client()))
	c.SetRequest(chronos.ChronosGetRepoProjectsRequest{})

	resp, err := c.GetRepoProjects()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Projects:", resp.Projects)
}
