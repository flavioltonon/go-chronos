package main

import (
	"flavioltonon/go-chronos/chronos"
	"flavioltonon/go-chronos/chronos/config/project"
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

	for _, project := range project.Projects() {
		c.SetRequest(chronos.ChronosGetProjectColumnsRequest{
			ProjectID: project.ID(),
		})

		resp, err := c.GetProjectColumns()
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Columns:", resp.Columns)
	}
}
