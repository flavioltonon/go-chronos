package chronos

import (
	"context"

	"github.com/flavioltonon/go-github/github"
)

func (chronos Chronos) ListLabels() ([]*github.Label, *github.Response, error) {
	return chronos.client.Repositories.ListLabels(context.Background(), GITHUB_REPOSITORY_OWNER, GITHUB_REPOSITORY_NAME, &github.ListOptions{})
}
