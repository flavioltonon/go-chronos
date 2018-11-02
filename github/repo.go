package github

type Repo struct {
	issue    Issue
	issues   []Issue
	newLabel Label
	labels   []Label
}
