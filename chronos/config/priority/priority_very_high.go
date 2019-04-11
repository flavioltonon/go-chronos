package priority

type PriorityVeryHigh struct{}

func init() {
	RegisterPriority(&PriorityVeryHigh{})
}

func (p PriorityVeryHigh) ID() int64 {
	return 1113351187
}

func (p PriorityVeryHigh) Name() string {
	return PRIORITY_LABEL_PRIORITY_VERY_HIGH
}

func (p PriorityVeryHigh) Deadline() Deadline {
	return Deadline{
		Duration:           24,
		Unit:               DEADLINE_TYPE_HOURS,
		DeduceNonWorkHours: DEDUCE_NON_WORK_HOURS_PRIORITY_VERY_HIGH,
	}
}
