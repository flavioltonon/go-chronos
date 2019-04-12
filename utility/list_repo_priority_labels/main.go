package main

import (
	"flavioltonon/go-chronos/chronos"
	"log"
	"os"
	"regexp"

	"github.com/flavioltonon/go-github/github"
)

func main() {
	var c chronos.Chronos

	auth := github.BasicAuthTransport{
		Username: os.Getenv("CHRONOS_GITHUB_LOGIN"),
		Password: os.Getenv("CHRONOS_GITHUB_PASSWORD"),
	}

	c.SetClient(github.NewClient(auth.Client()))

	labels, _, err := c.ListLabels()
	if err != nil {
		log.Println(err)
		return
	}

	var priorities = make([]chronos.Priority, 0)
	for _, label := range labels {
		if regexp.MustCompile(chronos.PRIORITY_LABEL_SIGNATURE).MatchString(label.GetName()) {
			priorities = append(priorities, chronos.Priority{
				ID:   label.GetID(),
				Name: label.GetName(),
			})
		}
	}
	log.Println("Priorities:", priorities)
}
