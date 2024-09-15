package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
	goWebAuthn "github.com/go-webauthn/webauthn/webauthn"

	"leoho.io/it16th-webauthn-rp-server/api"
	"leoho.io/it16th-webauthn-rp-server/database"
	"leoho.io/it16th-webauthn-rp-server/utils"
	"leoho.io/it16th-webauthn-rp-server/webauthn"
)

var assertionSessionData *goWebAuthn.SessionData

func StartAssertionHandler(ctx *gin.Context) {
	fmt.Println("call /assertion/options")

	var request *api.CredentialGetOptionsRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			api.CommonResponse{
				Status:       "failed",
				ErrorMessage: "failed to parse request body, error: " + err.Error(),
			},
		)
		return
	}
	fmt.Println("Parse request success")
	reqBody := utils.PrintJSON(request)
	fmt.Println("Request body: ", reqBody)

	foundUser, err := database.GetUserByName(request.Username)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			api.CommonResponse{
				Status:       "failed",
				ErrorMessage: "failed to get user by name, error: " + err.Error(),
			},
		)
		return
	}
	fmt.Println("get user by name success")

	options, sessionData, err := webauthn.WebAuthn.BeginLogin(foundUser)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			api.CommonResponse{
				Status:       "failed",
				ErrorMessage: "failed to begin login, error: " + err.Error(),
			},
		)
		return
	}
	fmt.Println("begin login success")

	assertionSessionData = sessionData

	err = database.UpdateUser(
		foundUser, &database.User{
			Challenge: options.Response.Challenge.String(),
		},
	)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			api.CommonResponse{
				Status:       "failed",
				ErrorMessage: "failed to update user, error: " + err.Error(),
			},
		)
		return
	}
	fmt.Println("update user success")

	ctx.JSON(
		http.StatusOK,
		api.CredentialGetOptionsResponse{
			CommonResponse: api.CommonResponse{
				Status:       "success",
				ErrorMessage: "",
			},
			PublicKeyCredentialRequestOptions: options.Response,
		},
	)
}

func FinishAssertionHandler(ctx *gin.Context) {
	fmt.Println("call /assertion/result")

	var request *api.AuthenticatorAssertionResponseRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			api.CommonResponse{
				Status:       "failed",
				ErrorMessage: "failed to parse request body, error: " + err.Error(),
			},
		)
		return
	}
	fmt.Println("Parse request success")
	reqBody := utils.PrintJSON(request)
	fmt.Println("Request body: ", reqBody)

	var authenticatorClientDataJSON []byte
	_, err := base64.RawURLEncoding.Decode(authenticatorClientDataJSON, request.Response.ClientDataJSON)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			api.CommonResponse{
				Status:       "failed",
				ErrorMessage: "failed to decode clientDataJSON, error: " + err.Error(),
			},
		)
		return
	}
	var clientDataJSON map[string]interface{}
	if err := json.Unmarshal(authenticatorClientDataJSON, &clientDataJSON); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			api.CommonResponse{
				Status:       "failed",
				ErrorMessage: "failed to unmarshal clientDataJSON, error: " + err.Error(),
			},
		)
		return
	}
	fmt.Println("Decode clientDataJSON success")

	if challenge, ok := clientDataJSON["challenge"].(string); !ok || challenge != attestationSessionData.Challenge {
		ctx.JSON(
			http.StatusBadRequest,
			api.CommonResponse{
				Status:       "failed",
				ErrorMessage: "challenge mismatch",
			},
		)
		return
	} else {
		foundUser, err := database.GetUserByChallenge(challenge)
		if err != nil {
			ctx.JSON(
				http.StatusInternalServerError,
				api.CommonResponse{
					Status:       "failed",
					ErrorMessage: "failed to get user by challenge, error: " + err.Error(),
				},
			)
			return
		}
		fmt.Println("get user by challenge success")

		car := protocol.CredentialAssertionResponse{
			PublicKeyCredential: protocol.PublicKeyCredential{
				Credential: protocol.Credential{
					ID:   request.Id,
					Type: request.Type,
				},
				RawID:                  []byte(request.Id),
				ClientExtensionResults: request.GetClientExtensionResults,
			},
			AssertionResponse: request.Response,
		}
		pca, err := car.Parse()
		if err != nil {
			ctx.JSON(
				http.StatusInternalServerError,
				api.CommonResponse{
					Status:       "failed",
					ErrorMessage: "failed to parse assertion response, error: " + err.Error(),
				},
			)
			return
		}
		fmt.Println("parse assertion response success")

		credential, err := webauthn.WebAuthn.ValidateLogin(foundUser, *assertionSessionData, pca)
		if err != nil {
			ctx.JSON(
				http.StatusInternalServerError,
				api.CommonResponse{
					Status:       "failed",
					ErrorMessage: "failed to validate login, error: " + err.Error(),
				},
			)
			return
		}
		fmt.Println("validate login success")

		credentialJSON := utils.PrintJSON(credential)
		fmt.Println("Credential: ", credentialJSON)

		ctx.JSON(
			http.StatusOK,
			api.CommonResponse{
				Status:       "success",
				ErrorMessage: "",
			},
		)
	}
}
