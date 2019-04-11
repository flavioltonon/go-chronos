package priority

type PriorityMedium struct{}

func init() {
	RegisterPriority(&PriorityMedium{})
}

func (p PriorityMedium) ID() int64 {
	return 1113350767
}

func (p PriorityMedium) Name() string {
	return PRIORITY_LABEL_PRIORITY_MEDIUM
}

func (p PriorityMedium) Deadline() Deadline {
	return Deadline{
		Duration:           15,
		Unit:               DEADLINE_TYPE_DAYS,
		DeduceNonWorkHours: DEDUCE_NON_WORK_HOURS_PRIORITY_MEDIUM,
	}
}
