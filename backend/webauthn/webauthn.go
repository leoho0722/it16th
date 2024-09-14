package webauthn

import (
	"fmt"

	"github.com/go-webauthn/webauthn/webauthn"

	"leoho.io/it16th-webauthn-rp-server/config"
)

var WebAuthn *webauthn.WebAuthn

func NewRPServer() {
	rpServer := config.GetWebAuthnConfiguration()
	c := &webauthn.Config{
		RPID:          rpServer.Id,
		RPDisplayName: rpServer.DisplayName,
		RPOrigins:     []string{rpServer.Origin},
	}

	webAuthn, err := webauthn.New(c)
	if err != nil {
		fmt.Println(err)
	}

	WebAuthn = webAuthn
}
