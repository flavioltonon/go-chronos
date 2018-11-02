package github

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/go-resty/resty"
)

func (h Repo) GetLabel(labelName string) error {
	u, _ := url.Parse(fmt.Sprintf("%s/repos/%s/%s/labels/", GITHUB_API_URL, OWNER, REPO))
	fullURL, _ := u.Parse(labelName)

	resp, err := resty.R().
		SetBasicAuth(os.Getenv("CHRONOS_GITHUB_LOGIN"), os.Getenv("CHRONOS_GITHUB_PASSWORD")).
		Get(fullURL.String())
	if err != nil {
		return ErrUnableToSendGetLabelRequest
	}

	var label Label
	err = json.Unmarshal(resp.Body(), &label)
	if err != nil {
		return ErrUnableToUnmarshalGetLabelResponse
	}

	log.Println(label)

	return nil
}
