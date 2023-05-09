package models

import whatsapp "github.com/sufficit/sufficit-quepasa/whatsapp"

type QpPictureResponse struct {
	QpResponse
	Info *whatsapp.WhatsappProfilePicture `json:"info,omitempty"`
}
