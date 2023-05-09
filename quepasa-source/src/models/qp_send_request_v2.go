package models

import (
	"net/http"

	whatsapp "github.com/sufficit/sufficit-quepasa/whatsapp"
)

type QpSendRequestV2 struct {
	Recipient  string         `json:"recipient,omitempty"`
	Message    string         `json:"message,omitempty"`
	Attachment QPAttachmentV1 `json:"attachment,omitempty"`
}

func (source *QpSendRequestV2) EnsureValidChatId(r *http.Request) (err error) {
	chatid, err := whatsapp.FormatEndpoint(source.Recipient)
	if err != nil {
		return
	}

	source.Recipient = chatid
	return
}

func (source *QpSendRequestV2) ToCurrentVersion() *QpSendRequest {
	result := &QpSendRequest{
		ChatId: source.Recipient,
		Text:   source.Message,
	}

	attachDefault := QPAttachmentV1{}
	if source.Attachment != attachDefault {
		// proccess attach
	}

	return result
}
