package github

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/go-resty/resty"
)

func (h Repo) GetIssues(query string) error {
	u, _ := url.Parse(fmt.Sprintf("%s/repos/%s/%s/issues", GITHUB_API_URL, OWNER, REPO))
	fullURL, _ := u.Parse(query)
	resp, err := resty.R().
		SetBasicAuth(os.Getenv("CHRONOS_GITHUB_LOGIN"), os.Getenv("CHRONOS_GITHUB_PASSWORD")).
		Get(fullURL.String())
	if err != nil {
		log.Println(err)
		return err
	}

	var issues []Issue
	err = json.Unmarshal(resp.Body(), &issues)
	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("%+v", issues))

	return nil
}
