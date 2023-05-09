package models

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	library "github.com/sufficit/sufficit-quepasa/library"
	whatsapp "github.com/sufficit/sufficit-quepasa/whatsapp"
)

// Serviço que controla os servidores / bots individuais do whatsapp
type QPWhatsappService struct {
	Servers     map[string]*QpWhatsappServer
	DB          *QpDatabase
	Initialized bool

	initlock   *sync.Mutex
	appendlock *sync.Mutex
}

var WhatsappService *QPWhatsappService

func QPWhatsappStart() error {
	if WhatsappService == nil {
		log.Trace("whatsapp service starting ...")

		db := GetDatabase()
		WhatsappService = &QPWhatsappService{
			Servers:    make(map[string]*QpWhatsappServer),
			DB:         db,
			initlock:   &sync.Mutex{},
			appendlock: &sync.Mutex{},
		}

		// seeding database
		err := InitialSeed()
		if err != nil {
			return err
		}
		// iniciando servidores e cada bot individualmente
		return WhatsappService.Initialize()
	} else {
		log.Debug("attempt to start whatsapp service, already started ...")
	}
	return nil
}

// Inclui um novo servidor em um serviço já em andamento
// *Usado quando se passa pela verificação do QRCode
// *Usado quando se inicializa o sistema
func (service *QPWhatsappService) AppendNewServer(info *QpServer) (server *QpWhatsappServer, err error) {
	server, ok := service.Servers[info.Token]
	if !ok {
		// Adiciona na lista de servidores
		log.Infof("adding new server on cache: %s, wid: %s", info.Token, info.WId)

		// Creating a new instance
		server, err = service.NewQpWhatsappServer(info)
		if err != nil {
			log.Errorf("error on append new server: %s, :: %s", info.WId, err.Error())
			return
		}

		service.Servers[info.Token] = server
	} else {
		// Adiciona na lista de servidores
		log.Infof("updating new server on cache: %s, wid: %s", info.Token, info.WId)

		server.QpServer = info
	}
	return
}

func (service *QPWhatsappService) AppendPaired(paired *QpWhatsappPairing) (server *QpWhatsappServer, err error) {
	server, ok := service.Servers[paired.Token]
	if !ok {
		// Adiciona na lista de servidores
		log.Infof("adding paired server on cache: %s, wid: %s", paired.Token, paired.WId)

		info := &QpServer{Token: paired.Token, WId: paired.WId}

		// Creating a new instance
		server, err = service.NewQpWhatsappServer(info)
		if err != nil {
			log.Errorf("error on append new server: %s, :: %s", info.WId, err.Error())
			return
		}

		service.Servers[info.Token] = server
	} else {
		// Adiciona na lista de servidores
		log.Infof("updating paired server on cache: %s, old wid: %s, new wid: %s", server.Token, server.WId, paired.WId)

		server.connection = paired.conn

		//server.Disconnect("old session")
		server.WId = paired.WId
	}

	server.Verified = true

	// checking user
	if paired.User != nil {
		server.User = paired.User.Username
	}

	err = server.Save()
	return
}

//region CONSTRUCTORS

// Instanciando um novo servidor para controle de whatsapp
func (service *QPWhatsappService) NewQpWhatsappServer(info *QpServer) (server *QpWhatsappServer, err error) {

	var serverLogLevel log.Level
	if info.Devel {
		serverLogLevel = log.DebugLevel
	} else {
		serverLogLevel = log.InfoLevel
	}

	serverLogger := log.New()
	serverLogger.SetLevel(serverLogLevel)
	serverLogEntry := serverLogger.WithField("token", info.Token)

	server = &QpWhatsappServer{
		QpServer:       info,
		syncConnection: &sync.Mutex{},
		syncMessages:   &sync.Mutex{},
		StartTime:      time.Now().UTC(),

		stopRequested: false,
		Log:           serverLogEntry,
		db:            service.DB.Servers,
	}

	server.HandlerEnsure()
	server.WebHookEnsure()
	server.WebhookFill(info.Token, service.DB.Webhooks)
	return
}

func (service *QPWhatsappService) GetOrCreateServerFromToken(token string) (server *QpWhatsappServer, err error) {
	log.Debugf("locating server: %s", token)
	server, ok := service.Servers[token]
	if !ok {
		log.Debugf("server: %s, not in cache, looking up database", token)
		exists, err := service.DB.Servers.Exists(token)
		if err != nil {
			return nil, err
		}

		var info *QpServer
		if exists {
			info, err = service.DB.Servers.FindByToken(token)
			if err != nil {
				return nil, err
			}
			log.Debugf("server: %s, found", token)
		} else {
			info = &QpServer{
				Token: token,
			}
		}

		server, err = service.AppendNewServer(info)
	}
	return
}

/*
<summary>

	Get or Create a server for scanned qrcode from forms with current user informations and a whatsapp section id
	* use same token if already exists

</summary>
*/
func (service *QPWhatsappService) GetOrCreateServer(user string, wid string) (result *QpWhatsappServer, err error) {
	log.Debugf("locating server with section id: %s", wid)

	phone := library.GetPhoneByWId(wid)
	log.Infof("wid to phone: %s", phone)

	var server *QpWhatsappServer
	servers := service.GetServersForUser(user)
	for _, item := range servers {
		if item.GetNumber() == phone {
			server = item
			server.WId = wid
			break
		}
	}

	if server == nil {
		token := uuid.New().String()
		log.Infof("creating new server with token: %s", token)
		info := &QpServer{
			Token:        token,
			User:         user,
			WId:          wid,
			HandleGroups: true,
		}

		server, err = service.AppendNewServer(info)
		if err != nil {
			return
		}
	}

	// server.Disconnect("GetOrCreateServer")
	result = server
	return
}

func (service *QPWhatsappService) Delete(server *QpWhatsappServer) (err error) {
	err = server.Delete()
	if err != nil {
		return
	}

	delete(service.Servers, server.Token)
	return
}

// Função que irá iniciar todos os servidores apartir do banco de dados
func (service *QPWhatsappService) Initialize() (err error) {

	if !service.Initialized {

		servers := service.DB.Servers.FindAll()
		for _, info := range servers {

			// appending server to cache
			server, err := service.AppendNewServer(info)
			if err != nil {
				return err
			}

			state := server.GetStatus()
			if state == whatsapp.UnPrepared || IsValidToStart(state) {

				// initialize individual server
				server.Log.Debugf("starting whatsapp server ... on %s state", state)
				go server.Initialize()
			} else {
				server.Log.Debugf("not auto starting cause state: %s", state)
			}
		}

		service.Initialized = true
	}

	return
}

// Função privada que irá iniciar todos os servidores apartir do banco de dados
func (service *QPWhatsappService) GetServersForUser(username string) (servers map[string]*QpWhatsappServer) {
	servers = make(map[string]*QpWhatsappServer)
	for _, server := range service.Servers {
		if server.GetOwnerID() == username {
			servers[server.Token] = server
		}
	}
	return
}

func (service *QPWhatsappService) FindByToken(token string) (*QpWhatsappServer, error) {
	for _, server := range service.Servers {
		if server.Token == token {
			return server, nil
		}
	}

	err := fmt.Errorf("server not found for token: %s", token)
	return nil, err
}

func (service *QPWhatsappService) GetUser(username string, password string) (user *QpUser, err error) {
	log.Debugf("finding user: %s", username)
	return service.DB.Users.Check(username, password)
}
