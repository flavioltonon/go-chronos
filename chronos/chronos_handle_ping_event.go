package chronos

import (
	"fmt"

	"github.com/flavioltonon/go-github/github"
)

func (chronos Chronos) HandlePingEvent(event interface{}) error {
	var pingEvent = event.(*github.PingEvent)
	fmt.Println("Event: Ping received.\nZen message:", pingEvent.GetZen())
	return nil
}
