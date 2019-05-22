package chronos

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/flavioltonon/go-github/github"
)

func (chronos *Chronos) ListenToGitHub(w http.ResponseWriter, r *http.Request) {
	var (
		event interface{}
		err   error
	)

	defer func() {
		defer r.Body.Close()
		if err != nil {
			log.Fatal(event, err)
		}
	}()

	payload, err := github.ValidatePayload(r, []byte(os.Getenv("GITHUB_WEBHOOK_SECRET")))
	if err != nil {
		err = fmt.Errorf("error validating request body: err=%s", err)
		return
	}

	event, err = github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		err = fmt.Errorf("could not parse webhook: err=%s", err)
		return
	}

	auth := github.BasicAuthTransport{
		Username: os.Getenv("GITHUB_LOGIN"),
		Password: os.Getenv("GITHUB_PASSWORD"),
	}

	chronos.SetClient(github.NewClient(auth.Client()))

	if len(holidays) == 0 {
		chronos.SetRequest(ChronosGetHolidaysRequest{
			Country: os.Getenv("HOLIDAY_STANDARD_COUNTRY"),
			Year:    strconv.Itoa(time.Now().Local().Year()),
		})

		resp, err := chronos.GetHolidays()
		if err != nil {
			log.Println("failed to get holidays")
		}

		holidays = resp.Holidays
	}

	switch event.(type) {
	case *github.IssuesEvent:
		err = chronos.HandleIssuesEvent(event.(*github.IssuesEvent))
	case *github.LabelEvent:
		err = chronos.HandleLabelEvent(event.(*github.LabelEvent))
	case *github.ProjectCardEvent:
		err = chronos.HandleProjectCardEvent(event.(*github.ProjectCardEvent))
	case *github.PingEvent:
		err = chronos.HandlePingEvent(event.(*github.PingEvent))
	default:
		err = fmt.Errorf("unknown event type %s", github.WebHookType(r))
	}
}
