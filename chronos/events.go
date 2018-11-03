package chronos

type Event interface {
	HandleEvent()
}

type IssuesEventHandler struct{}

type LabelEventHandler struct{}

type ProjectCardEventHandler struct{}

type PingEventHandler struct{}
