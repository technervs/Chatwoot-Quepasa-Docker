package controllers

import (
	"html/template"
	"net/http"
	"strings"

	websocket "github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	models "github.com/sufficit/sufficit-quepasa/models"
)

func RedirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, FormLoginEndpoint, http.StatusFound)
}

// Google chrome bloqueou wss, portanto retornaremos sempre ws apatir de agora
func WebSocketProtocol() string {
	protocol := "ws"
	isSecure := models.ENV.UseSSLForWebSocket()
	if isSecure {
		protocol = "wss"
	}

	return protocol
}

//
// Cycle
//

// CycleHandler renders route POST "/bot/cycle"
func FormCycleController(w http.ResponseWriter, r *http.Request) {
	_, server, err := GetUserAndServer(w, r)
	if err != nil {
		// retorno já tratado pela funcao
		return
	}

	err = server.CycleToken()
	if err != nil {
		RespondServerError(server, w, err)
		return
	}

	http.Redirect(w, r, FormAccountEndpoint, http.StatusFound)
}

// DebugHandler renders route POST "/bot/debug"
func FormDebugController(w http.ResponseWriter, r *http.Request) {
	_, server, err := GetUserAndServer(w, r)
	if err != nil {
		// retorno já tratado pela funcao
		return
	}

	_, err = server.ToggleDevel()
	if err != nil {
		RespondServerError(server, w, err)
		return
	}

	http.Redirect(w, r, FormAccountEndpoint, http.StatusFound)
}

// ToggleHandler renders route POST "/bot/toggle"
func FormToggleController(w http.ResponseWriter, r *http.Request) {
	_, server, err := GetUserAndServer(w, r)
	if err != nil {
		// retorno já tratado pela funcao
		return
	}

	err = server.Toggle()
	if err != nil {
		response := &models.QpResponse{}
		response.ParseError(err)
		RespondInterface(w, response)
		return
	}

	http.Redirect(w, r, FormAccountEndpoint, http.StatusFound)
}

// ToggleHandler renders route POST "/bot/toggle"
func FormToggleBroadcastController(w http.ResponseWriter, r *http.Request) {
	_, server, err := GetUserAndServer(w, r)
	if err != nil {
		// retorno já tratado pela funcao
		return
	}

	_, err = server.ToggleBroadcast()
	if err != nil {
		RespondServerError(server, w, err)
		return
	}

	http.Redirect(w, r, FormAccountEndpoint, http.StatusFound)
}

func FormToggleGroupsController(w http.ResponseWriter, r *http.Request) {
	_, server, err := GetUserAndServer(w, r)
	if err != nil {
		// retorno já tratado pela funcao
		return
	}

	_, err = server.ToggleGroups()
	if err != nil {
		RespondServerError(server, w, err)
		return
	}

	http.Redirect(w, r, FormAccountEndpoint, http.StatusFound)
}

//
// Verify
//

// VerifyFormHandler renders route GET "/bot/verify" ?mode={sd|md}
func VerifyFormHandler(w http.ResponseWriter, r *http.Request) {
	data := models.QPFormVerifyData{
		PageTitle:   "Verify To Add or Update",
		Protocol:    WebSocketProtocol(),
		Host:        r.Host,
		Destination: FormAccountEndpoint,
	}

	templates := template.Must(template.ParseFiles(
		"views/layouts/main.tmpl",
		"views/bot/verify.tmpl",
	))
	templates.ExecuteTemplate(w, "main", data)
}

// VerifyHandler renders route GET "/bot/verify/ws"
func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	user, err := models.GetFormUser(r)
	if err != nil {
		log.Errorf("connection upgrade error (not logged): %s", err.Error())
		RedirectToLogin(w, r)
		return
	}

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("(websocket): service error: %s", err.Error())
		return
	}

	WebSocketStart(user, conn)
}

//
// Delete
//

// DeleteHandler renders route POST "/form/delete"
func FormDeleteController(w http.ResponseWriter, r *http.Request) {
	_, server, err := GetUserAndServer(w, r)
	if err != nil {
		// retorno já tratado pela funcao
		return
	}

	server.Log.Warnf("delete requested by form !")
	if err := models.WhatsappService.Delete(server); err != nil {
		RespondServerError(server, w, err)
		return
	}

	http.Redirect(w, r, FormAccountEndpoint, http.StatusFound)
}

//
// Helpers
//

// Facilitador que traz usuario e servidor para quem esta autenticado
func GetUserAndServer(w http.ResponseWriter, r *http.Request) (user *models.QpUser, server *models.QpWhatsappServer, err error) {
	user, err = models.GetFormUser(r)
	if err != nil {
		RedirectToLogin(w, r)
		return
	}

	r.ParseForm()

	token := GetToken(r)
	server, err = models.WhatsappService.FindByToken(token)
	if err != nil {
		return
	}

	return
}

func GetServerFromRequest(r *http.Request) (server *models.QpWhatsappServer, err error) {
	token := GetToken(r)
	return models.WhatsappService.FindByToken(token)
}

func GetDownloadPrefix(token string) (path string) {
	path = "/download?token={token}&cache=false&messageid={messageid}"
	path = strings.Replace(path, "{token}", token, -1)
	path = strings.Replace(path, "{messageid}", "", -1)
	return
}
