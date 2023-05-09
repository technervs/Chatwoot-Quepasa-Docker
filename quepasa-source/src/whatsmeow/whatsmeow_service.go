package whatsmeow

import (
	"context"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"

	library "github.com/sufficit/sufficit-quepasa/library"
	whatsapp "github.com/sufficit/sufficit-quepasa/whatsapp"
	whatsmeow "go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type WhatsmeowServiceModel struct {
	Container *sqlstore.Container
}

var WhatsmeowService *WhatsmeowServiceModel

func (service *WhatsmeowServiceModel) Start() {
	if service == nil {
		log.Info("Starting Whatsmeow Service ....")

		dbLog := waLog.Stdout("whatsmeow/database", string(WarnLevel), true)

		// check if exists old whatsmeow.db
		var cs string
		if _, err := os.Stat("whatsmeow.db"); err == nil {
			cs = "file:whatsmeow.db?_foreign_keys=on"
		} else {
			// using new quepasa.sqlite
			cs = "file:whatsmeow.sqlite?_foreign_keys=on"
		}

		container, err := sqlstore.New("sqlite3", cs, dbLog)
		if err != nil {
			panic(err)
		}

		WhatsmeowService = &WhatsmeowServiceModel{Container: container}

		showing := whatsapp.WhatsappWebAppName
		if len(whatsapp.WhatsappWebAppSystem) > 0 {
			showing += " " + whatsapp.WhatsappWebAppSystem
		}

		var version [3]uint32
		version[0] = 0
		version[1] = 9
		version[2] = 0
		store.SetOSInfo(showing, version)
	}
}

// Used for scan QR Codes
// Dont forget to attach handlers after success login
func (service *WhatsmeowServiceModel) CreateEmptyConnection() (conn *WhatsmeowConnection, err error) {
	logger := log.StandardLogger()
	logger.SetLevel(log.DebugLevel)

	return service.CreateConnection("", logger.WithContext(context.Background()))
}

func (service *WhatsmeowServiceModel) CreateConnection(wid string, loggerEntry *log.Entry) (conn *WhatsmeowConnection, err error) {
	clientLog := waLog.Stdout("Whatsmeow/Client", loggerEntry.Level.String(), true)
	client, err := service.GetWhatsAppClient(wid, clientLog)
	if err != nil {
		return
	}

	if len(wid) > 0 {
		loggerEntry = loggerEntry.WithField("wid", wid)
	}

	handlers := &WhatsmeowHandlers{
		Client: client,
		log:    loggerEntry,
	}

	err = handlers.Register()
	if err != nil {
		return
	}

	conn = &WhatsmeowConnection{
		Client:   client,
		Handlers: handlers,
		waLogger: clientLog,
		log:      loggerEntry,
	}

	client.PrePairCallback = conn.PairedCallBack
	return
}

func (service *WhatsmeowServiceModel) GetStoreFromWid(wid string) (str *store.Device, err error) {
	if wid == "" {
		str = service.Container.NewDevice()
	} else {
		devices, err := service.Container.GetAllDevices()
		if err != nil {
			err = fmt.Errorf("(Whatsmeow)(EX001) error on getting store from wid (%s): %v", wid, err)
			return str, err
		} else {
			for _, element := range devices {
				if element.ID.String() == wid {
					str = element
					break
				}
			}

			if str == nil {
				err = &WhatsmeowStoreNotFoundException{Wid: wid}
				return str, err
			}
		}
	}

	return
}

// Temporaly
func (service *WhatsmeowServiceModel) GetStoreForMigrated(phone string) (str *store.Device, err error) {

	devices, err := service.Container.GetAllDevices()
	if err != nil {
		err = fmt.Errorf("(Whatsmeow)(EX001) error on getting store from phone (%s): %v", phone, err)
		return str, err
	} else {
		for _, element := range devices {
			if library.GetPhoneByWId(element.ID.String()) == phone {
				str = element
				break
			}
		}

		if str == nil {
			err = &WhatsmeowStoreNotFoundException{Wid: phone}
			return str, err
		}
	}

	return
}

func (service *WhatsmeowServiceModel) GetWhatsAppClient(wid string, logger waLog.Logger) (client *whatsmeow.Client, err error) {
	deviceStore, err := service.GetStoreFromWid(wid)
	if deviceStore != nil {
		client = whatsmeow.NewClient(deviceStore, logger)
		client.AutoTrustIdentity = true
		client.EnableAutoReconnect = true
	}
	return
}

// Flush entire Whatsmeow Database
// Use with wisdom !
func (service *WhatsmeowServiceModel) FlushDatabase() (err error) {
	devices, err := service.Container.GetAllDevices()
	if err != nil {
		return
	}

	for _, element := range devices {
		err = element.Delete()
		if err != nil {
			return
		}
	}

	return
}
