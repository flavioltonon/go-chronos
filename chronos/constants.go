package chronos

const (
	GITHUB_API_URL = "https://api.github.com"
	OWNER          = "flavioltonon"
	REPO           = "go-chronos"

	STANDARD_TIME_LOCATION = "America/Sao_Paulo"

	WORK_HOURS_INITIAL = 9
	WORK_HOURS_FINAL   = 18

	PRIORITY_LABEL_SIGNATURE             = "Prioridade"
	PRIORITY_LABEL_PRIORIDADE_MUITO_ALTA = "Prioridade: Muito Alta"
	PRIORITY_LABEL_PRIORIDADE_ALTA       = "Prioridade: Alta"
	PRIORITY_LABEL_PRIORIDADE_MEDIA      = "Prioridade: Média"
	PRIORITY_LABEL_PRIORIDADE_BAIXA      = "Prioridade: Baixa"

	DEADLINE_TYPE_HOURS = "horas"
	DEADLINE_TYPE_DAYS  = "dias"

	DEADLINE_LABEL_SIGNATURE             = "Prazo"
	DEADLINE_LABEL_PRIORIDADE_MUITO_ALTA = "24 horas"
	DEADLINE_LABEL_PRIORIDADE_ALTA       = "3 dias"
	DEADLINE_LABEL_PRIORIDADE_MEDIA      = "15 dias"
	DEADLINE_LABEL_PRIORIDADE_BAIXA      = "60 dias"
	DEADLINE_LABEL_OVERDUE               = "Overdue"

	COLUMN_BACKLOG       = "Backlog"
	COLUMN_SPRINTBACKLOG = "Sprint backlog"
	COLUMN_DEPLOY        = "Deploy"
	COLUMN_DONE          = "Concluído"

	STATUS_LABEL_SIGNATURE = "Status"
	STATUS_LABEL_DEPLOY    = "Status: PR/Testes"

	HOLIDAY_API_URL = "https://holidayapi.pl/v1/holidays"
)
