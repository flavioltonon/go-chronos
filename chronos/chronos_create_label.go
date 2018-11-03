package chronos

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-resty/resty"
)

func (h *Chronos) CreateLabel(name string) error {
	var labelSpec = LabelSpec{
		Name:  name,
		Color: SetColorToLabel(name),
	}

	resp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody(labelSpec).
		SetBasicAuth(os.Getenv("CHRONOS_GITHUB_LOGIN"), os.Getenv("CHRONOS_GITHUB_PASSWORD")).
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
