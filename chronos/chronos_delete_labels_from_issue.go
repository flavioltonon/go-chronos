package chronos

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-resty/resty"
)

func (h Chronos) DeleteLabelsFromIssue(number int, labelNames []string) error {
	for _, labelName := range labelNames {
		resp, err := resty.R().
			SetHeader("Content-Type", "application/json").
			SetBasicAuth(os.Getenv("CHRONOS_GITHUB_LOGIN"), os.Getenv("CHRONOS_GITHUB_PASSWORD")).
			Delete(fmt.Sprintf("%s/repos/%s/%s/issues/%s/labels/%s", GITHUB_API_URL, OWNER, REPO, strconv.Itoa(number), labelName))
		if err != nil {
			return ErrUnableToSendDeleteLabelsFromIssueRequest
		}

		if resp.IsError() {
			return ErrDeleteLabelsFromIssueBadResponse
		}
	}

	return nil
}
