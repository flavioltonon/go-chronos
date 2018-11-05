package chronos

import (
	"fmt"

	"github.com/flavioltonon/go-github/github"
)

func (h *PingEventHandler) HandleEvent(event interface{}) error {
	pingEvent := event.(*github.PingEvent)
	fmt.Println("Event: Ping received.\nZen message:", pingEvent.GetZen())
	return nil
}
