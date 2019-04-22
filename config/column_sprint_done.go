package config

type ColumnSprintDone struct{}

func init() {
	RegisterColumn(&ColumnSprintDone{})
}

func (p ColumnSprintDone) ID() int64 {
	return 5039444
}

func (p ColumnSprintDone) Name() string {
	return "Sprint Done"
}

func (p ColumnSprintDone) Project() int64 {
	return 1908642
}

func (p ColumnSprintDone) StandardIssueState() string {
	return "closed"
}
