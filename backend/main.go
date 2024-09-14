package main

import (
	"leoho.io/it16th-webauthn-rp-server/config"
	"leoho.io/it16th-webauthn-rp-server/database"
	"leoho.io/it16th-webauthn-rp-server/route"
	"leoho.io/it16th-webauthn-rp-server/webauthn"
)

func main() {
	config.Parse()
	database.Connect()
	webauthn.NewRPServer()
	route.NewRoute()
}
