package route

import (
	"github.com/gin-gonic/gin"

	"leoho.io/it16th-webauthn-rp-server/controller"
)

func NewRoute() {
	app := gin.Default()

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

	app.Run("0.0.0.0:8080")
}
