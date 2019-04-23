package chronos

type PriorityNone struct{}

func init() {
	RegisterPriority(&PriorityNone{})
}

func (p PriorityNone) ID() int64 {
	return 1332172722
}

func (p PriorityNone) Name() string {
	return "Priority: None"
}

func (p PriorityNone) Level() int {
	return 99
}

func (p PriorityNone) Deadline() Deadline {
	return Deadline{}
}
