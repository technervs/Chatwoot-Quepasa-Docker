package main

import (
	"github.com/joho/godotenv"
	controllers "github.com/sufficit/sufficit-quepasa/controllers"
	models "github.com/sufficit/sufficit-quepasa/models"
	whatsapp "github.com/sufficit/sufficit-quepasa/whatsapp"
	whatsmeow "github.com/sufficit/sufficit-quepasa/whatsmeow"

	log "github.com/sirupsen/logrus"
)

// @title chi-swagger example APIs
// @version 1.0
// @description chi-swagger example APIs
// @BasePath /
func main() {

	// Carregando variaveis de ambiente apartir de arquivo .env
	godotenv.Load()

	if models.ENV.DEBUGJsonMessages() {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	// Verifica se é necessario realizar alguma migração de base de dados
	err := models.MigrateToLatest()
	if err != nil {
		log.Fatalf("Database migration error: %s", err.Error())
	}

	// should became before whatsmeow start
	title := models.ENV.AppTitle()
	if len(title) > 0 {
		whatsapp.WhatsappWebAppSystem = title
	}

	whatsmeow.WhatsmeowService.Start()

	// must execute after whatsmeow started
	for _, element := range models.Running {
		if handler, ok := models.MigrationHandlers[element]; ok {
			handler(element)
		}
	}

	// Inicializando serviço de controle do whatsapp
	// De forma assíncrona
	err = models.QPWhatsappStart()
	if err != nil {
		log.Fatalf("Whatsapp service starting error: %s", err.Error())
	}

	controllers.QPWebServerStart()
	log.Info("Ready !")
}
