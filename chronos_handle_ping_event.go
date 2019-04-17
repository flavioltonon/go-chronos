package chronos

import (
	"fmt"

	"github.com/flavioltonon/go-github/github"
)

func (chronos Chronos) HandlePingEvent(event *github.PingEvent) error {
	fmt.Println("Event: Ping received.\nZen message:", event.GetZen())

	return nil
}
