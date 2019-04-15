package priority

type PriorityHigh struct{}

func init() {
	RegisterPriority(&PriorityHigh{})
}

func (p PriorityHigh) ID() int64 {
	return 1113351048
}

func (p PriorityHigh) Name() string {
	return "Priority: High"
}

func (p PriorityHigh) Level() int {
	return 1
}

func (p PriorityHigh) Deadline() Deadline {
	return Deadline{
		Duration: 3,
		Unit:     "days",
	}
}
