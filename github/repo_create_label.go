package github

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty"
)

func (h *Repo) CreateLabel(name string) error {
	var labelSpec = LabelSpec{
		Name:  name,
		Color: SetColorToLabel(name),
	}

	resp, err := resty.R().
		SetBody(labelSpec).
		SetBasicAuth(CHRONOS_GITHUB_LOGIN, CHRONOS_GITHUB_PASSWORD).
		Post(fmt.Sprintf("%s/repos/%s/%s/labels", GITHUB_API_URL, OWNER, REPO))
	if err != nil {
		return ErrUnableToSendCreateLabelRequest
	}

	var label Label
	err = json.Unmarshal(resp.Body(), &label)
	if err != nil {
		return ErrUnableToUnmarshalCreateLabelResponse
	}

	h.newLabel = label

	return nil
}
