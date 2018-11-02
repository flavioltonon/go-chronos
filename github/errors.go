package github

import "errors"

var (
	ErrUnableToSendGetLabelsFromIssueRequest       = errors.New("unable to send GetLabelsFromIssue request")
	ErrUnableToUnmarshalGetLabelsFromIssueResponse = errors.New("unable to unmarshal GetLabelsFromIssue response")
	ErrUnableToSendAddLabelsToIssueRequest         = errors.New("unable to send AddLabelsToIssue request")
	ErrUnableToUnmarshalAddLabelsToIssueResponse   = errors.New("unable to unmarshal AddLabelsToIssue response")
	ErrUnableToSendDeleteLabelsFromIssueRequest    = errors.New("unable to send DeleteLabelsFromIssue request")
	ErrUnableToDefineTimer                         = errors.New("unable to define timer from unexpected label")
	ErrUnableToSendGetIssueRequest                 = errors.New("unable to send GetIssue request")
	ErrUnableToUnmarshalGetIssueResponse           = errors.New("unable to unmarshal GetIssue response")
	ErrUnableToSendGetLabelRequest                 = errors.New("unable to send GetLabel request")
	ErrUnableToUnmarshalGetLabelResponse           = errors.New("unable to unmarshal GetLabel response")
	ErrUnableToSendCreateLabelRequest              = errors.New("unable to send CreateLabel request")
	ErrUnableToUnmarshalCreateLabelResponse        = errors.New("unable to unmarshal CreateLabel response")
)
