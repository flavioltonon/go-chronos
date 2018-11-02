package github

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-resty/resty"
)

func (h Repo) GetLabels() error {
	resp, err := resty.R().
		SetBasicAuth(os.Getenv("CHRONOS_GITHUB_LOGIN"), os.Getenv("CHRONOS_GITHUB_PASSWORD")).
		Get(fmt.Sprintf("%s/repos/%s/%s/labels", GITHUB_API_URL, OWNER, REPO))
	if err != nil {
		log.Println(err)
		return err
	}

	var labels []Label
	err = json.Unmarshal(resp.Body(), &labels)
	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("%+v", labels))

	return nil
}
