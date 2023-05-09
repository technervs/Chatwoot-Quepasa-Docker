package models

import (
	"encoding/base64"
)

type QpSendRequestEncoded struct {
	QpSendRequest
	Content string `json:"content"`
}

func (source *QpSendRequestEncoded) GenerateContent() (err error) {
	content, err := base64.StdEncoding.DecodeString(source.Content)
	if err != nil {
		return
	}

	source.QpSendRequest.Content = content
	return
}
