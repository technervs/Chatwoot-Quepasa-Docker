package whatsmeow

import (
	"fmt"
)

type WhatsmeowStoreNotFoundException struct {
	Wid string
}

func (e *WhatsmeowStoreNotFoundException) Error() string {
	return fmt.Sprint("cant find a store")
}

func (e *WhatsmeowStoreNotFoundException) Unauthorized() bool {
	return true
}
