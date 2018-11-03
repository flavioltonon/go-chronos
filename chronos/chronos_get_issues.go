package chronos

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-resty/resty"
)

func (h *Chronos) GetIssues(query map[string]string) error {
	resp, err := resty.R().
		SetQueryParams(query).
		SetBasicAuth(os.Getenv("CHRONOS_GITHUB_LOGIN"), os.Getenv("CHRONOS_GITHUB_PASSWORD")).
		Get(fmt.Sprintf("%s/repos/%s/%s/issues", GITHUB_API_URL, OWNER, REPO))
	if err != nil {
		return ErrUnableToSendGetIssuesRequest
	}

	var issues []Issue
	err = json.Unmarshal(resp.Body(), &issues)
	if err != nil {
		return ErrUnableToUnmarshalGetIssuesResponse
	}

	h.issues = issues

	return nil
}
