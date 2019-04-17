package priority

type PriorityLow struct{}

func init() {
	RegisterPriority(&PriorityLow{})
}

func (p PriorityLow) ID() int64 {
	return 1113350430
}

func (p PriorityLow) Name() string {
	return "Priority: Low"
}

func (p PriorityLow) Level() int {
	return 3
}

func (p PriorityLow) Deadline() Deadline {
	return Deadline{
		Duration: 60,
		Unit:     "days",
	}
}
