package chronos

import (
	"os"
	"strconv"

	"github.com/flavioltonon/go-github/github"
)

type Chronos struct {
	auth   github.BasicAuthTransport
	client *github.Client

	request  interface{}
	response interface{}
}

func (h *Chronos) SetClient(client *github.Client) {
	h.client = client
}

func (h *Chronos) SetRequest(request interface{}) {
	h.request = request
}

func (h Chronos) Response() interface{} {
	return h.response
}

func (h Chronos) UserID() int64 {
	user, err := strconv.ParseInt(os.Getenv("CHRONOS_GITHUB_USER_ID"), 10, 64)
	if err != nil {
		panic("invalid CHRONOS_GITHUB_USER_ID")
	}
	return user
}
