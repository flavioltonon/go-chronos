package chronos

type ProjectAcessoOnboarding struct{}

func init() {
	RegisterProject(&ProjectAcessoOnboarding{})
}

func (p ProjectAcessoOnboarding) ID() int64 {
	return 1302676
}

func (p ProjectAcessoOnboarding) Name() string {
	return "Acesso Onboarding"
}
