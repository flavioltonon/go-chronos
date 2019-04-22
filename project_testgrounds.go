package chronos

type ProjectTestgrounds struct{}

func init() {
	RegisterProject(&ProjectTestgrounds{})
}

func (p ProjectTestgrounds) ID() int64 {
	return 1908642
}

func (p ProjectTestgrounds) Name() string {
	return "Testgrounds"
}
