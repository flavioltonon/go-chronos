package github

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-resty/resty"
)

func (h *Repo) GetLabelsFromIssue(number string) error {
	resp, err := resty.R().
		SetBasicAuth(CHRONOS_GITHUB_LOGIN, CHRONOS_GITHUB_PASSWORD).
		Get(fmt.Sprintf("%s/repos/%s/%s/issues/%s/labels", GITHUB_API_URL, OWNER, REPO, number))
	if err != nil {
		return ErrUnableToSendGetLabelsFromIssueRequest
	}

	var labels []Label
	err = json.Unmarshal(resp.Body(), &labels)
	if err != nil {
		return ErrUnableToUnmarshalGetLabelsFromIssueResponse
	}

	h.labels = labels

	log.Println(labels)

	return nil
}