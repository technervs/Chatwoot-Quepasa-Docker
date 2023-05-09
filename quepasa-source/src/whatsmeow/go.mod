module github.com/sufficit/sufficit-quepasa/whatsmeow

require (
	github.com/sufficit/sufficit-quepasa/library v0.0.0-00010101000000-000000000000
	github.com/sufficit/sufficit-quepasa/whatsapp v0.0.0-00010101000000-000000000000
)

require (
	filippo.io/edwards25519 v1.0.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/gosimple/slug v1.13.1
	github.com/gosimple/unidecode v1.0.1 // indirect
	github.com/mattn/go-sqlite3 v1.14.16
	github.com/sirupsen/logrus v1.9.0
	go.mau.fi/libsignal v0.1.0 // indirect
	go.mau.fi/whatsmeow v0.0.0-20230407182255-e4dca20d3923
	golang.org/x/crypto v0.7.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/protobuf v1.30.0
)

replace github.com/sufficit/sufficit-quepasa/whatsmeow => ./

replace github.com/sufficit/sufficit-quepasa/whatsapp => ../whatsapp

replace github.com/sufficit/sufficit-quepasa/library => ../library

go 1.19