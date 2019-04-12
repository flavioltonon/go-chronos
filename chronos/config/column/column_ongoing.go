package column

type ColumnOngoing struct{}

func init() {
	RegisterColumn(&ColumnOngoing{})
}

func (p ColumnOngoing) ID() int64 {
	return 0
}

func (p ColumnOngoing) Name() string {
	return "Ongoing"
}

func (p ColumnOngoing) StandardIssueState() string {
	return "open"
}
