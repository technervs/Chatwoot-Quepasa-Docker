package whatsapp

// Eventos vindos do serviço de whatsapp
type IWhatsappHandlers interface {

	// Recebimento/Envio de mensagem
	Message(*WhatsappMessage)

	// Event
	LoggedOut(string)
}
