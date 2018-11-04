package chronos

import "errors"

var (
	ErrNothingToUpdate = errors.New("nothing to update")

	ErrUnableToSendGetLabelsFromIssueRequest       = errors.New("unable to send GetLabelsFromIssue request")
	ErrUnableToUnmarshalGetLabelsFromIssueResponse = errors.New("unable to unmarshal GetLabelsFromIssue response")

	ErrUnableToAddLabelsToIssue                  = errors.New("unable to add labels to issue")
	ErrUnableToSendAddLabelsToIssueRequest       = errors.New("unable to send AddLabelsToIssue request")
	ErrUnableToUnmarshalAddLabelsToIssueResponse = errors.New("unable to unmarshal AddLabelsToIssue response")
	ErrAddLabelsToIssueBadResponse               = errors.New("got bad response during AddLabelsToIssue request")

	ErrUnableToDeleteLabelsFromIssue            = errors.New("unable to delete labels from issue")
	ErrUnableToSendDeleteLabelsFromIssueRequest = errors.New("unable to send DeleteLabelsFromIssue request")
	ErrDeleteLabelsFromIssueBadResponse         = errors.New("got bad response during DeleteLabelsFromIssue request")

	ErrUnableToDefineTimer = errors.New("unable to define timer from unexpected label")

	ErrUnableToGetIssue                  = errors.New("unable to get issue")
	ErrUnableToUnmarshalGetIssueResponse = errors.New("unable to unmarshal GetIssue response")

	ErrUnableToGetIssuesFromRepo          = errors.New("unable to get issues from repository")
	ErrUnableToUnmarshalGetIssuesResponse = errors.New("unable to unmarshal GetIssues response")

	ErrUnableToSendGetLabelRequest       = errors.New("unable to send GetLabel request")
	ErrUnableToUnmarshalGetLabelResponse = errors.New("unable to unmarshal GetLabel response")

	ErrUnableToSendCreateLabelRequest       = errors.New("unable to send CreateLabel request")
	ErrUnableToUnmarshalCreateLabelResponse = errors.New("unable to unmarshal CreateLabel response")

	ErrUnableToSendGetHolidaysRequest       = errors.New("unable to send GetHolidays request")
	ErrUnableToUnmarshalGetHolidaysResponse = errors.New("unable to unmarshal GetHolidays response")

	ErrUnableToGetRepositoryProjects = errors.New("unable to get repository projects")
	ErrUnableToGetProjectCards       = errors.New("unable to get project cards")
	ErrUnableToGetProjectColumns     = errors.New("unable to get project columns")
	ErrUnexpectedProjectColumnName   = errors.New("unexpected project column name")

	ErrUnableToUpdateIssueState = errors.New("unable to update issue state")
)
