package chronos

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-resty/resty"
)

func (h Chronos) AddLabelsToIssue(number int, labelNames []string) error {
	resp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody(labelNames).
		SetBasicAuth(os.Getenv("CHRONOS_GITHUB_LOGIN"), os.Getenv("CHRONOS_GITHUB_PASSWORD")).
		Post(fmt.Sprintf("%s/repos/%s/%s/issues/%s/labels", GITHUB_API_URL, OWNER, REPO, strconv.Itoa(number)))
	if err != nil {
		return ErrUnableToSendAddLabelsToIssueRequest
	}

	if resp.IsError() {
		return ErrAddLabelsToIssueBadResponse
	}

	return nil
}
