module github.com/sufficit/sufficit-quepasa/models

require (
	github.com/go-chi/chi/v5 v5.0.7
	github.com/joncalhoun/migrate v0.0.2
	github.com/sufficit/sufficit-quepasa/library v0.0.0-00010101000000-000000000000
	github.com/sufficit/sufficit-quepasa/whatsapp v0.0.0-00010101000000-000000000000
	github.com/sufficit/sufficit-quepasa/whatsmeow v0.0.0-00010101000000-000000000000
)

require (
	github.com/go-chi/chi v1.5.4 // indirect
	github.com/gosimple/slug v1.13.1 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
)

require (
	filippo.io/edwards25519 v1.0.0 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/go-chi/jwtauth v4.0.4+incompatible
	github.com/go-sql-driver/mysql v1.5.0 // indirect
	github.com/google/uuid v1.1.1
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.5.2
	github.com/mattn/go-sqlite3 v2.0.3+incompatible
	github.com/sirupsen/logrus v1.9.0
	github.com/skip2/go-qrcode v0.0.0-20191027152451-9434209cb086
	go.mau.fi/libsignal v0.1.0 // indirect
	go.mau.fi/whatsmeow v0.0.0-20230407182255-e4dca20d3923 // indirect
	golang.org/x/crypto v0.7.0
	golang.org/x/sys v0.7.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

replace github.com/sufficit/sufficit-quepasa/library => ../library

replace github.com/sufficit/sufficit-quepasa/whatsmeow => ../whatsmeow

replace github.com/sufficit/sufficit-quepasa/whatsapp => ../whatsapp

replace github.com/sufficit/sufficit-quepasa/models => ./

go 1.19