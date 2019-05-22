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

	projects   map[int64]Project
	columns    map[int64]Column
	priorities map[int64]Priority
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
	user, err := strconv.ParseInt(os.Getenv("GITHUB_USER_ID"), 10, 64)
	if err != nil {
		panic("invalid GITHUB_USER_ID")
	}
	return user
}

func (h Chronos) Projects() map[int64]Project {
	return h.projects
}

func (h Chronos) Columns() map[int64]Column {
	return h.columns
}

func (h Chronos) Priorities() map[int64]Priority {
	return h.priorities
}
