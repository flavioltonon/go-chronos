package priority

type PriorityLow struct{}

func init() {
	RegisterPriority(&PriorityLow{})
}

func (p PriorityLow) ID() int64 {
	return 1113350430
}

func (p PriorityLow) Name() string {
	return PRIORITY_LABEL_PRIORITY_LOW
}

func (p PriorityLow) Deadline() Deadline {
	return Deadline{
		Duration:           60,
		Unit:               DEADLINE_TYPE_DAYS,
		DeduceNonWorkHours: DEDUCE_NON_WORK_HOURS_PRIORITY_LOW,
	}
}
