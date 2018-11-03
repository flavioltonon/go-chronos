package main

import (
	"flavioltonon/go-chronos/chronos"
	"log"
	"os"

	"github.com/google/go-github/github"
)

func main() {
	var chronos chronos.Chronos

	auth := github.BasicAuthTransport{
		Username: os.Getenv("CHRONOS_GITHUB_LOGIN"),
		Password: os.Getenv("CHRONOS_GITHUB_PASSWORD"),
	}
	chronos.SetClient(github.NewClient(auth.Client()))

	err := chronos.UpdateIssuesDeadlines()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("All deadlines have been updated successfully")
}
