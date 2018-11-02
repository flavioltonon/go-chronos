package github

import (
	"fmt"

	"github.com/go-resty/resty"
)

func (h Repo) DeleteLabelsFromIssue(number string, labelNames []string) error {
	for _, labelName := range labelNames {
		_, err := resty.R().
			SetBasicAuth(CHRONOS_GITHUB_LOGIN, CHRONOS_GITHUB_PASSWORD).
			Delete(fmt.Sprintf("%s/repos/%s/%s/issues/%s/labels/%s", GITHUB_API_URL, OWNER, REPO, number, labelName))
		if err != nil {
			return ErrUnableToSendDeleteLabelsFromIssueRequest
		}
	}
	return nil
}
