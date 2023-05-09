package models

type QPSendRequest struct {
	Recipient  string       `json:"recipient,omitempty"`
	Message    string       `json:"message,omitempty"`
	Attachment QPAttachment `json:"attachment,omitempty"`
}
