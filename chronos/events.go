package chronos

type Event interface {
	HandleEvent()
}

type IssuesEventHandler struct{}

type ProjectCardEventHandler struct{}

type PingEventHandler struct{}
