package models

import (
	. "github.com/sufficit/sufficit-quepasa/whatsapp"
)

// Mensagem no formato QuePasa
// Utilizada na API do QuePasa para troca com outros sistemas
type QPAttachment struct {
	WhatsappAttachment

	// Public URL to direct download without encryption
	DirectPath string `json:"url,omitempty"`
}
