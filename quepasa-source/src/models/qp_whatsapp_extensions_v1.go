package models

import (
	"encoding/base64"

	whatsapp "github.com/sufficit/sufficit-quepasa/whatsapp"
)

func ToWhatsappAttachment(source *QPAttachmentV1) (attach *whatsapp.WhatsappAttachment, err error) {
	attach = &whatsapp.WhatsappAttachment{}
	content, err := base64.StdEncoding.DecodeString(source.Base64)
	if err != nil {
		return
	}

	attach.SetContent(&content)
	attach.FileName = source.FileName
	attach.Mimetype = source.MIME
	attach.FileLength = uint64(len(content))
	return
}

func ToWhatsappMessageV1(source *QpSendRequestV2) (msg *whatsapp.WhatsappMessage, err error) {
	recipient, err := whatsapp.FormatEndpoint(source.Recipient)
	if err != nil {
		return
	}

	attach, err := ToWhatsappAttachment(&source.Attachment)
	if err != nil {
		return
	}

	chat := whatsapp.WhatsappChat{Id: recipient}
	msg = &whatsapp.WhatsappMessage{}
	msg.Text = source.Message
	msg.Chat = chat
	if attach != nil {
		msg.Attachment = attach
	}
	return
}
