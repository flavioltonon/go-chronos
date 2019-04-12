package chronos

import (
	"context"
	"os"

	"github.com/flavioltonon/go-github/github"
)

func (chronos Chronos) ListLabels() ([]*github.Label, *github.Response, error) {
	return chronos.client.Repositories.ListLabels(
		context.Background(),
		os.Getenv("GITHUB_REPOSITORY_OWNER"),
		os.Getenv("GITHUB_REPOSITORY_NAME"),
		&github.ListOptions{},
	)
}
