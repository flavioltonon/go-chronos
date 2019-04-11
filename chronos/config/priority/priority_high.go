package priority

type PriorityHigh struct{}

func init() {
	RegisterPriority(&PriorityHigh{})
}

func (p PriorityHigh) ID() int64 {
	return 1113351048
}

func (p PriorityHigh) Name() string {
	return PRIORITY_LABEL_PRIORITY_HIGH
}

func (p PriorityHigh) Deadline() Deadline {
	return Deadline{
		Duration:           3,
		Unit:               DEADLINE_TYPE_DAYS,
		DeduceNonWorkHours: DEDUCE_NON_WORK_HOURS_PRIORITY_HIGH,
	}
}
