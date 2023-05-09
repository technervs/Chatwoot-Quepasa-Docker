package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	models "github.com/sufficit/sufficit-quepasa/models"
)

type errorResponse struct {
	Result string `json:"result"`
}

func RespondBadRequest(w http.ResponseWriter, err error) {
	log.Infof("Request Bad Format: %s", err.Error())
	RespondErrorCode(w, err, http.StatusBadRequest)
}

func RespondUnauthorized(w http.ResponseWriter, err error) {
	log.Debugf("Request Unauthorized: %s", err.Error())
	RespondErrorCode(w, err, http.StatusUnauthorized)
}

func RespondNoContent(w http.ResponseWriter, err error) {
	log.Debugf("Request Not found: %s", err.Error())
	RespondErrorCode(w, err, http.StatusNoContent)
}

/// Usado para avisar que o bot ainda não esta pronto
func RespondNotReady(w http.ResponseWriter, err error) {
	RespondErrorCode(w, err, http.StatusServiceUnavailable)
}

func RespondServerError(server *models.QpWhatsappServer, w http.ResponseWriter, err error) {
	if strings.Contains(err.Error(), "invalid websocket") {

		// Desconexão forçado é algum evento iniciado pelo whatsapp
		log.Errorf("(%s) Desconexão forçada por motivo de websocket inválido ou sem resposta", server.Token)
		go server.Restart()

	} else {
		if models.ENV.DEBUGRequests() {
			log.Errorf("(%s) !Request Server error: %s", server.Token, err)
		}
	}
	RespondErrorCode(w, err, http.StatusInternalServerError)
}

func RespondErrorCode(w http.ResponseWriter, err error, code int) {
	res := &errorResponse{
		Result: err.Error(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(res)
}

/*
<summary>
	Default response method
	Used in v3 *models.QpResponse
	Returns OK | Bad Request
</summary>
*/
func RespondInterfaceCode(w http.ResponseWriter, response interface{}, code uint) {

	// setting default reponse type as json
	w.Header().Set("Content-Type", "application/json")

	// Writing header code
	if code == 0 {
		code = http.StatusOK
		if qpresponse, ok := response.(models.QpResponseBasicInterface); ok {
			if !qpresponse.IsSuccess() {
				code = http.StatusBadRequest
			}
		}
	}

	w.WriteHeader(int(code))
	json.NewEncoder(w).Encode(response)
}

func RespondInterface(w http.ResponseWriter, response interface{}) {
	RespondInterfaceCode(w, response, 0)
}

func RespondSuccess(w http.ResponseWriter, response interface{}) {
	if qpresponse, ok := response.(models.QpResponseInterface); ok {
		if !qpresponse.IsSuccess() {
			if len(qpresponse.GetStatusMessage()) == 0 {
				qpresponse.ParseSuccess("")
			}
		}
	}

	RespondInterfaceCode(w, response, http.StatusOK)
}
