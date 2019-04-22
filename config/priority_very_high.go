package config

type PriorityVeryHigh struct{}

func init() {
	RegisterPriority(&PriorityVeryHigh{})
}

func (p PriorityVeryHigh) ID() int64 {
	return 1113351187
}

func (p PriorityVeryHigh) Name() string {
	return "Priority: Very High"
}

func (p PriorityVeryHigh) Level() int {
	return 0
}

func (p PriorityVeryHigh) Deadline() Deadline {
	return Deadline{
		Duration: 24,
		Unit:     "hours",
	}
}
