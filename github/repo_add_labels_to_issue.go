package github

import (
	"fmt"

	"github.com/go-resty/resty"
)

func (h Repo) AddLabelsToIssue(number string, labelNames []string) error {
	_, err := resty.R().
		SetBody(labelNames).
		SetBasicAuth(CHRONOS_GITHUB_LOGIN, CHRONOS_GITHUB_PASSWORD).
		Post(fmt.Sprintf("%s/repos/%s/%s/issues/%s/labels", GITHUB_API_URL, OWNER, REPO, number))
	if err != nil {
		return ErrUnableToSendAddLabelsToIssueRequest
	}

	return nil
}
