package controllers

import (
	"fmt"

	whatsapp "github.com/sufficit/sufficit-quepasa/whatsapp"
)

type ApiServerNotReadyException struct {
	Wid    string
	Status whatsapp.WhatsappConnectionState
}

func (e *ApiServerNotReadyException) Error() string {
	return fmt.Sprintf("server (%s) not ready yet ! current status: %s.", e.Wid, e.Status.String())
}
