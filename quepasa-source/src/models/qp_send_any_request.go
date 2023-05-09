package models

import (
	"encoding/base64"
	"io"
	"net/http"
	"path"
)

/*
<summary>
	Request to send any type of message
	1º Attachment Url
	2º Attachment Content Base64
	3º Text Message
</summary>
*/
type QpSendAnyRequest struct {
	QpSendRequest
	Url     string `json:"url,omitempty"`
	Content string `json:"content,omitempty"`
}

func (source *QpSendAnyRequest) GenerateEmbbedContent() (err error) {
	content, err := base64.StdEncoding.DecodeString(source.Content)
	if err != nil {
		return
	}

	source.QpSendRequest.Content = content
	return
}

func (source *QpSendAnyRequest) GenerateUrlContent() (err error) {
	resp, err := http.Get(source.Url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	source.QpSendRequest.Content = content

	// setting filename if empty
	if len(source.QpSendRequest.FileName) == 0 {
		source.QpSendRequest.FileName = path.Base(source.Url)
	}

	return
}
