package chronos

import "github.com/google/go-github/github"

type Chronos struct {
	auth   github.BasicAuthTransport
	client *github.Client

	request interface{}

	holidays Holidays

	issue  *github.Issue
	issues []*github.Issue

	newLabel *github.Label
	labels   []*github.Label

	timer   string
	overdue bool
}

func (h *Chronos) SetClient(client *github.Client) {
	h.client = client
}
