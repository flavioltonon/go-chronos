package chronos

type ColumnOngoing struct{}

func init() {
	RegisterColumn(&ColumnOngoing{})
}

func (p ColumnOngoing) ID() int64 {
	return 2244699
}

func (p ColumnOngoing) Name() string {
	return "Ongoing"
}

func (p ColumnOngoing) Project() int64 {
	return 1302676
}

func (p ColumnOngoing) StandardIssueState() string {
	return "open"
}
