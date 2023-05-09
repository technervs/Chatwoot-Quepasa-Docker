package models

import (
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	library "github.com/sufficit/sufficit-quepasa/library"
	whatsapp "github.com/sufficit/sufficit-quepasa/whatsapp"
)

type QpSendRequest struct {
	// (Optional) Used if passed
	Id string `json:"id,omitempty"`

	// Recipient of this message
	ChatId string `json:"chatId"`

	// (Optional) TrackId - less priority (urlparam -> query -> header -> body)
	TrackId string `json:"trackId,omitempty"`

	Text string `json:"text,omitempty"`

	// (Optional) Sugested filename on user download
	FileName string `json:"fileName,omitempty"`

	Content []byte
}

func (source *QpSendRequest) EnsureChatId(r *http.Request) (err error) {
	if len(source.ChatId) == 0 {
		source.ChatId = GetChatId(r)
	}

	if len(source.ChatId) == 0 {
		err = fmt.Errorf("chat id missing")
	}
	return
}

func (source *QpSendRequest) EnsureValidChatId(r *http.Request) (err error) {
	err = source.EnsureChatId(r)
	if err != nil {
		return
	}

	chatid, err := whatsapp.FormatEndpoint(source.ChatId)
	if err != nil {
		return
	}

	source.ChatId = chatid
	return
}

func (source *QpSendRequest) ToWhatsappMessage() (msg *whatsapp.WhatsappMessage, err error) {
	chatId, err := whatsapp.FormatEndpoint(source.ChatId)
	if err != nil {
		return
	}

	chat := whatsapp.WhatsappChat{Id: chatId}
	msg = &whatsapp.WhatsappMessage{
		Id:           source.Id,
		TrackId:      source.TrackId,
		Text:         source.Text,
		Chat:         chat,
		FromMe:       true,
		FromInternal: true,
	}

	// setting default type
	if len(msg.Text) > 0 {
		msg.Type = whatsapp.TextMessageType
	}

	return
}

func (source *QpSendRequest) ToWhatsappAttachment() (attach *whatsapp.WhatsappAttachment, err error) {
	attach = &whatsapp.WhatsappAttachment{}

	mimeType := library.GetMimeTypeFromContent(source.Content, source.FileName)

	// adjusting codec for ptt audio messages
	if strings.Contains(mimeType, "ogg") && !strings.Contains(mimeType, "opus") {
		mimeType = "audio/ogg; codecs=opus"
	}

	log.Tracef("detected mime type: %s, filename: %s", mimeType, source.FileName)
	filename := source.FileName

	// Defining a filename if not found before
	if len(filename) == 0 {
		filename = library.GenerateFileNameFromMimeType(mimeType)
	}

	attach.FileName = filename
	attach.FileLength = uint64(len(source.Content))
	attach.Mimetype = mimeType
	attach.SetContent(&source.Content)
	return
}
