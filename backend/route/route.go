package route

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"leoho.io/it16th-webauthn-rp-server/config"
	"leoho.io/it16th-webauthn-rp-server/controller"
)

func NewRoute() {
	app := gin.Default()

	wellknown := app.Group("/.well-known")
	{
		wellknown.GET("/apple-app-site-association", controller.AppleWellKnownHandler)
		wellknown.GET("/assetlinks.json", controller.AndroidWellKnownHandler)
	}

	attestation := app.Group("/attestation")
	{
		attestation.POST("/options", controller.StartAttestationHandler)
		attestation.POST("/result", controller.FinishAttestationHandler)
	}

	assertion := app.Group("/assertion")
	{
		assertion.POST("/options", controller.StartAssertionHandler)
		assertion.POST("/result", controller.FinishAssertionHandler)
	}

	serverConfig := config.GetServerConfiguration()
	addr := fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port)
	app.Run(addr)
}
