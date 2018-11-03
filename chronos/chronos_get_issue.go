package chronos

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-resty/resty"
)

func (h *Chronos) GetIssue(number string) error {
	resp, err := resty.R().
		SetBasicAuth(os.Getenv("CHRONOS_GITHUB_LOGIN"), os.Getenv("CHRONOS_GITHUB_PASSWORD")).
		Get(fmt.Sprintf("%s/repos/%s/%s/issues/%s", GITHUB_API_URL, OWNER, REPO, number))
	if err != nil {
		return ErrUnableToSendGetIssueRequest
	}

	var issue Issue
	err = json.Unmarshal(resp.Body(), &issue)
	if err != nil {
		return ErrUnableToUnmarshalGetIssueResponse
	}

	h.issue = issue

	return nil
}
