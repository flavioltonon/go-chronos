package github

import "time"

type Issue struct {
	ID      int       `json:"id"`
	Number  int       `json:"number"`
	Labels  []Label   `json:"labels"`
	Created time.Time `json:"created_at"`
}
