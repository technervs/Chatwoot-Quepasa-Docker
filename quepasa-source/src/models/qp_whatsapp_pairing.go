package models

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	library "github.com/sufficit/sufficit-quepasa/library"
	whatsapp "github.com/sufficit/sufficit-quepasa/whatsapp"
)

type QpWhatsappPairing struct {
	// Public token
	Token string `db:"token" json:"token" validate:"max=100"`

	// Whatsapp session id
	WId string `db:"wid" json:"wid" validate:"max=255"`

	User *QpUser `json:"user,omitempty"`

	conn whatsapp.IWhatsappConnection `json:"-"`
}

func (source *QpWhatsappPairing) OnPaired(wid string) {
	source.WId = wid

	// updating token if from user
	if source.User != nil {
		source.Token = source.GetUserToken()
	}

	log.Infof("paired whatsapp section %s, for token %s", source.WId, source.Token)
	server, err := WhatsappService.AppendPaired(source)
	if err != nil {
		log.Errorf("paired error: %s", err.Error())
		return
	}

	go server.EnsureReady()
}

func (source *QpWhatsappPairing) GetConnection() (whatsapp.IWhatsappConnection, error) {
	if source.conn == nil {
		conn, err := NewEmptyConnection(source.OnPaired)
		if err != nil {
			return nil, err
		}
		source.conn = conn
	}

	return source.conn, nil
}

func (source *QpWhatsappPairing) GetUserToken() string {
	phone := library.GetPhoneByWId(source.WId)
	log.Infof("wid to phone: %s", phone)

	servers := WhatsappService.GetServersForUser(source.User.Username)
	for _, item := range servers {
		if item.GetNumber() == phone {
			return item.Token
		}
	}

	return uuid.New().String()
}
