package chronos

import "github.com/google/go-github/github"

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

func (h Chronos) UserID() int64 {
	return CHRONOS_GITHUB_USER_ID
}
