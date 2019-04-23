package chronos

type PriorityMedium struct{}

func init() {
	RegisterPriority(&PriorityMedium{})
}

func (p PriorityMedium) ID() int64 {
	return 843114598
}

func (p PriorityMedium) Name() string {
	return "Priority: Medium"
}

func (p PriorityMedium) Level() int {
	return 2
}

func (p PriorityMedium) Deadline() Deadline {
	return Deadline{
		Duration: 15,
		Unit:     "days",
	}
}
