package chronos

type Priority struct {
	ID       int64
	Name     string
	Deadline Deadline
}

type Deadline struct {
	Duration int
	Unit     string
}

// Custom label and rules settings
func (chronos Chronos) Priorities() []Priority {
	return []Priority{
		{
			ID:   1113351187,
			Name: PRIORITY_LABEL_PRIORITY_VERY_HIGH,
			Deadline: Deadline{
				Duration: 24,
				Unit:     DEADLINE_TYPE_HOURS,
			},
		},
		{
			ID:   1113351048,
			Name: PRIORITY_LABEL_PRIORITY_HIGH,
			Deadline: Deadline{
				Duration: 3,
				Unit:     DEADLINE_TYPE_DAYS,
			},
		},
		{
			ID:   1113350767,
			Name: PRIORITY_LABEL_PRIORITY_MEDIUM,
			Deadline: Deadline{
				Duration: 15,
				Unit:     DEADLINE_TYPE_DAYS,
			},
		},
		{
			ID:   1113350430,
			Name: PRIORITY_LABEL_PRIORITY_LOW,
			Deadline: Deadline{
				Duration: 60,
				Unit:     DEADLINE_TYPE_DAYS,
			},
		},
	}
}
