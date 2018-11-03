package chronos

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/go-resty/resty"
)

func (h *Chronos) GetLabelsFromIssue(number int) error {
	resp, err := resty.R().
		SetBasicAuth(os.Getenv("CHRONOS_GITHUB_LOGIN"), os.Getenv("CHRONOS_GITHUB_PASSWORD")).
		Get(fmt.Sprintf("%s/repos/%s/%s/issues/%s/labels", GITHUB_API_URL, OWNER, REPO, strconv.Itoa(number)))
	if err != nil {
		return ErrUnableToSendGetLabelsFromIssueRequest
	}

	var labels []Label
	err = json.Unmarshal(resp.Body(), &labels)
	if err != nil {
		return ErrUnableToUnmarshalGetLabelsFromIssueResponse
	}

	h.labels = labels

	return nil
}
