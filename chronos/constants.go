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
	DEADLINE_LABEL_PRIORIDADE_MUITO_ALTA = "Prazo: 24 horas"
	DEADLINE_LABEL_PRIORIDADE_ALTA       = "Prazo: 3 dias"
	DEADLINE_LABEL_PRIORIDADE_MEDIA      = "Prazo: 15 dias"
	DEADLINE_LABEL_PRIORIDADE_BAIXA      = "Prazo: 60 dias"
	DEADLINE_LABEL_OVERDUE               = "Overdue"

	DEDUCE_NON_WORK_HOURS_PRIORIDADE_MUITO_ALTA = false
	DEDUCE_NON_WORK_HOURS_PRIORIDADE_ALTA       = true
	DEDUCE_NON_WORK_HOURS_PRIORIDADE_MEDIA      = true
	DEDUCE_NON_WORK_HOURS_PRIORIDADE_BAIXA      = true

	COLUMN_BACKLOG        = "Backlog"
	COLUMN_SPRINT_BACKLOG = "Sprint backlog"
	COLUMN_DEPLOY         = "Deploy"
	COLUMN_DONE           = "Concluído"

	STANDARD_ISSUE_STATE_COLUMN_BACKLOG         = "open"
	STANDARD_ISSUE_STATE_COLUMN__SPRINT_BACKLOG = "open"
	STANDARD_ISSUE_STATE_COLUMN_DEPLOY          = "open"
	STANDARD_ISSUE_STATE_COLUMN_DONE            = "closed"

	STATUS_LABEL_SIGNATURE      = "Status"
	STATUS_LABEL_BACKLOG        = ""
	STATUS_LABEL_SPRINT_BACKLOG = ""
	STATUS_LABEL_DEPLOY         = "Status: PR/Testes"
	STATUS_LABEL_DONE           = ""

	HOLIDAY_API_URL = "https://holidayapi.pl/v1/holidays"
)
