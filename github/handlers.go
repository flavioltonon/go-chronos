package github

type Chronos struct{}

type Event interface {
	HandleEvent()
}

type IssuesEventHandler struct{}

type ProjectCardEventHandler struct{}

type PingEventHandler struct{}
