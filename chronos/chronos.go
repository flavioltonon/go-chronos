package chronos

type Chronos struct {
	holidays Holidays

	issue  Issue
	issues []Issue

	newLabel Label
	labels   []Label

	timer   string
	overdue bool
}
