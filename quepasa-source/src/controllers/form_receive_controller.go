package controllers

import (
	"html/template"
	"net/http"

	. "github.com/sufficit/sufficit-quepasa/metrics"
	models "github.com/sufficit/sufficit-quepasa/models"
	whatsapp "github.com/sufficit/sufficit-quepasa/whatsapp"
)

// FormReceiveController renders route GET "/bot/{token}/receive"
func FormReceiveController(w http.ResponseWriter, r *http.Request) {
	data := models.QPFormReceiveData{PageTitle: "Receive", FormAccountEndpoint: FormAccountEndpoint}

	server, err := GetServerFromRequest(r)
	if err != nil {
		data.ErrorMessage = err.Error()
	} else {
		data.Number = server.GetWid()
		data.Token = server.Token
		data.DownloadPrefix = GetDownloadPrefix(server.Token)
	}

	// Evitando tentativa de download de anexos sem o bot estar devidamente sincronizado
	status := server.GetStatus()
	if status != whatsapp.Ready {
		RespondNotReady(w, &ApiServerNotReadyException{Wid: server.GetWid(), Status: status})
		return
	}

	queryValues := r.URL.Query()
	paramTimestamp := queryValues.Get("timestamp")
	timestamp, err := GetTimestamp(paramTimestamp)
	if err != nil {
		MessageReceiveErrors.Inc()
		RespondServerError(server, w, err)
		return
	}

	messages := GetMessages(server, timestamp)
	data.Messages = messages

	MessagesReceived.Add(float64(len(messages)))

	templates := template.Must(template.ParseFiles(
		"views/layouts/main.tmpl",
		"views/bot/receive.tmpl"))
	templates.ExecuteTemplate(w, "main", data)
}
