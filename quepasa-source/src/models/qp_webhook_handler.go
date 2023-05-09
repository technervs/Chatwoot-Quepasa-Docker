package models

import (
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
	whatsapp "github.com/sufficit/sufficit-quepasa/whatsapp"
)

type QPWebhookHandler struct {
	server *QpWhatsappServer
}

func (w *QPWebhookHandler) Handle(payload *whatsapp.WhatsappMessage) {
	if !w.HasWebhook() {
		return
	}

	if payload.Type == whatsapp.DiscardMessageType|whatsapp.UnknownMessageType {
		log.Debugf("ignoring unknown message type on webhook request: %v", reflect.TypeOf(&payload))
		return
	}

	if payload.Type == whatsapp.TextMessageType && len(strings.TrimSpace(payload.Text)) <= 0 {
		log.Debug("ignoring empty text message on webhook request: %v", payload.Id)
		return
	}

	if payload.Chat.Id == "status@broadcast" && !w.server.HandleBroadcast {
		log.Debug("ignoring broadcast message on webhook request: %v", payload.Id)
		return
	}

	PostToWebHookFromServer(w.server, payload)
}

func (w *QPWebhookHandler) HasWebhook() bool {
	if w.server != nil {
		return len(w.server.Webhooks) > 0
	}
	return false
}
