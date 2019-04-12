# go-chronos

## Running
- Download all dependencies: ```dep ensure```
- Set all necessary environment variables, then run ```source configmap && source env```

### Variables:
- CHRONOS_GITHUB_USER_ID: can be found by logging any kind of event made by the account designed to be used by Chronos
- CHRONOS_GITHUB_LOGIN: Chronos's GitHub username/e-mail
- CHRONOS_GITHUB_PASSWORD: Chronos's GitHub password
- GITHUB_WEBHOOK_SECRET: webhook secret that has been configured in the settings of the project that will be handled by Chronos

### NGROK
- Run ngrok at port 8080: ```ngrok.exe http 8080```
PS: NGROK_STANDARD_PORT can be changed at the application configmap

### Event listener - The GitHub event interpreter
- Run the main file of the event listening module: ```go run ./evlistener/main.go```

**Standard implemented functions:**
- If a human user sets a priority label for an issue, Chronos will associate a deadline label to it, based on the standard deadline of that priority;
- If a human user adds a deadline label from an issue, Chronos will remove it;
- If a human user removes a deadline label from an issue, Chronos will relabel it;
- If a human user closes an issue, Chronos will reopen it;
- If a human user reopens an issue, Chronos will close it again;
- If Chronos reopens an issue, it will recalculate its deadline;
- If a human user moves an issue from a project column to another, it will have its status updated - will be opened or closed based on the new column STANDARD_ISSUE_STATE.

### Job dispenser - The issues deadlines updater
- Run the main file of the job dispenser module: ```go run ./jobdispenser/main.go```

**Standard implemented functions:**
- Every single open issue inside the configured project will have its deadline updated based on its creation date.
