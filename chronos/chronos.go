package chronos

import "github.com/flavioltonon/go-github/github"

type Chronos struct {
	auth   github.BasicAuthTransport
	client *github.Client

	request interface{}
}

func (h *Chronos) SetClient(client *github.Client) {
	h.client = client
}

func (h *Chronos) SetRequest(request interface{}) {
	h.request = request
}
