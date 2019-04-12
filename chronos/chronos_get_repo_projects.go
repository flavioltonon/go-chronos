package chronos

import (
	"context"
	"os"

	"github.com/flavioltonon/go-github/github"
)

type ChronosGetRepoProjectsRequest struct {
	ProjectID int64

	client *github.Client

	projects map[int64]string

	issue *github.Issue

	issueState       string
	issueStatusLabel string
}

type ChronosGetRepoProjectsResponse struct {
	Projects map[int64]string `json:"projects"`
}

func (h *ChronosGetRepoProjectsRequest) getRepoProjects() error {
	repoProjects, _, err := h.client.Repositories.ListProjects(
		context.Background(),
		os.Getenv("GITHUB_REPOSITORY_OWNER"),
		os.Getenv("GITHUB_REPOSITORY_NAME"),
		nil,
	)
	if err != nil {
		return ErrUnableToGetRepoProjects
	}

	projects := make(map[int64]string)
	for _, project := range repoProjects {
		projects[project.GetID()] = project.GetName()
	}

	h.projects = projects

	return nil
}

func (h Chronos) GetRepoProjects() (resp ChronosGetRepoProjectsResponse, err error) {
	var req = h.request.(ChronosGetRepoProjectsRequest)

	req.client = h.client

	err = req.getRepoProjects()
	resp.Projects = req.projects

	return resp, err
}
