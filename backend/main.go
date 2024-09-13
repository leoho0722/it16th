package main

import (
	"leoho.io/it16th-webauthn-rp-server/database"
	"leoho.io/it16th-webauthn-rp-server/route"
)

func main() {
	database.Connect()
	route.NewRoute()
}
