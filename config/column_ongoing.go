package config

type ColumnOngoing struct{}

func init() {
	RegisterColumn(&ColumnOngoing{})
}

func (p ColumnOngoing) ID() int64 {
	return 5073004
}

func (p ColumnOngoing) Name() string {
	return "Ongoing"
}

func (p ColumnOngoing) Project() int64 {
	return 1908642
}

func (p ColumnOngoing) StandardIssueState() string {
	return "open"
}
