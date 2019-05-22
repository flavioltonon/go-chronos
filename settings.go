package chronos

type Setting struct {
	Kind string
	Data config
}

type config interface{}

type Settings []Setting

type Project struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Column struct {
	ID                 int64  `json:"id"`
	Name               string `json:"name"`
	Project            int64  `json:"project"`
	StandardIssueState string `json:"standardIssueState"`
}

type Priority struct {
	ID       int64    `json:"id"`
	Name     string   `json:"name"`
	Level    int      `json:"level"`
	Deadline Deadline `json:"deadline"`
}

type Deadline struct {
	Duration  int    `json:"duration"`
	Unit      string `json:"unit"`
	Updatable bool   `json:"updatable"`
}
