package controllers

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	. "github.com/sufficit/sufficit-quepasa/metrics"
	. "github.com/sufficit/sufficit-quepasa/models"
	. "github.com/sufficit/sufficit-quepasa/whatsapp"
)

func FormSendController(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		controllerHttpGet(w, r)
	case "POST":
		controllerHttpPost(w, r)
	}
}

// Renders route GET "/bot/{token}/send"
func controllerHttpGet(w http.ResponseWriter, r *http.Request) {
	data := QPFormSendData{PageTitle: "Send"}

	server, err := GetServerFromRequest(r)
	if err != nil {
		data.ErrorMessage = err.Error()
		renderSendForm(w, data)
		return
	} else {
		data.Server = server.QpServer
	}

	renderSendForm(w, data)
}

// Renders route POST "/bot/{token}/send"
// Vindo do formulário de testes
func controllerHttpPost(w http.ResponseWriter, r *http.Request) {
	data := QPFormSendData{PageTitle: "Send"}

	server, err := GetServerFromRequest(r)
	if err != nil {
		data.ErrorMessage = err.Error()
		renderSendForm(w, data)
		return
	} else {
		data.Server = server.QpServer
	}

	attachment, err := uploadFile(w, r)
	if err != nil {
		data.ErrorMessage = err.Error()
		renderSendForm(w, data)
		return
	}

	r.ParseForm()
	recipient := r.Form.Get("recipient")
	message := r.Form.Get("message")

	msg, err := ToWhatsappMessage(recipient, message, attachment)
	if err != nil {
		RespondServerError(server, w, err)
		return
	}

	_, err = server.SendMessage(msg)
	if err != nil {
		RespondServerError(server, w, err)
		return
	}

	data.MessageId = msg.GetId()

	// Increment counter statistics
	MessagesSent.Inc()

	renderSendForm(w, data)
}

func uploadFile(w http.ResponseWriter, r *http.Request) (attach *WhatsappAttachment, err error) {
	log.Trace("form post, checking for file")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		return
	}

	// FormFile returns the first file for the given key `attachment`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, reader, err := r.FormFile("attachment")
	if err != nil {
		if strings.Contains(err.Error(), "no such file") {
			err = nil
			return
		}
		return
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	attach = &WhatsappAttachment{}
	attach.SetContent(&content)
	attach.Mimetype = reader.Header.Get("content-type")
	attach.FileLength = uint64(reader.Size)
	attach.FileName = reader.Filename
	return
}

func renderSendForm(w http.ResponseWriter, data QPFormSendData) {
	templates := template.Must(template.ParseFiles("views/layouts/main.tmpl", "views/bot/send.tmpl"))
	templates.ExecuteTemplate(w, "main", data)
}
