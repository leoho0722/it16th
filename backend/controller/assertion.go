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

	authenticationOptions := func(options *protocol.PublicKeyCredentialRequestOptions) {
		options.UserVerification = protocol.UserVerificationRequirement(request.UserVerification)
	}
	options, sessionData, err := webauthn.WebAuthn.BeginLogin(foundUser, authenticationOptions)
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

	sessionData.Challenge = base64.RawStdEncoding.EncodeToString([]byte(sessionData.Challenge))
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

	fmt.Println("Response: ", utils.PrintJSON(options.Response))

	ctx.JSON(
		http.StatusOK,
		api.CredentialGetOptionsResponse{
			CommonResponse: api.CommonResponse{
				Status:       "ok",
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

	fmt.Println("AssertionSessionData: ", utils.PrintJSON(assertionSessionData))

	authenticatorData, err := base64.RawURLEncoding.DecodeString(request.Response.AuthenticatorData)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			api.CommonResponse{
				Status:       "failed",
				ErrorMessage: "failed to decode authenticatorData, error: " + err.Error(),
			},
		)
		return
	}
	fmt.Println("Decode authenticatorData success")

	authenticatorSignature, err := base64.RawURLEncoding.DecodeString(request.Response.Signature)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			api.CommonResponse{
				Status:       "failed",
				ErrorMessage: "failed to decode signature, error: " + err.Error(),
			},
		)
		return
	}
	fmt.Println("Decode signature success")

	authenticatorUserHandle, err := base64.RawURLEncoding.DecodeString(request.Response.UserHandle)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			api.CommonResponse{
				Status:       "failed",
				ErrorMessage: "failed to decode userHandle, error: " + err.Error(),
			},
		)
		return
	}
	fmt.Println("Decode userHandle success")

	authenticatorClientDataJSON, err := base64.RawURLEncoding.DecodeString(request.Response.ClientDataJSON)
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

	challenge, ok := clientDataJSON["challenge"].(string)
	if !ok {
		ctx.JSON(
			http.StatusBadRequest,
			api.CommonResponse{
				Status:       "failed",
				ErrorMessage: "challenge not found",
			},
		)
		return
	}

	if challenge != assertionSessionData.Challenge {
		ctx.JSON(
			http.StatusBadRequest,
			api.CommonResponse{
				Status:       "failed",
				ErrorMessage: "challenge mismatch",
			},
		)
		return
	} else {
		decodedChallenge, err := base64.RawURLEncoding.DecodeString(challenge)
		if err != nil {
			ctx.JSON(
				http.StatusBadRequest,
				api.CommonResponse{
					Status:       "failed",
					ErrorMessage: "failed to decode challenge, error: " + err.Error(),
				},
			)
			return
		}
		challenge = string(decodedChallenge)

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

		credentialRawID, err := utils.DecodeToBase64StdEncoding(request.Id)
		if err != nil {
			ctx.JSON(
				http.StatusBadRequest,
				api.CommonResponse{
					Status:       "failed",
					ErrorMessage: "failed to encode credential raw id, error: " + err.Error(),
				},
			)
			return
		}

		car := protocol.CredentialAssertionResponse{
			PublicKeyCredential: protocol.PublicKeyCredential{
				Credential: protocol.Credential{
					ID:   request.Id,
					Type: request.Type,
				},
				RawID:                  protocol.URLEncodedBase64(credentialRawID),
				ClientExtensionResults: request.GetClientExtensionResults,
			},
			AssertionResponse: protocol.AuthenticatorAssertionResponse{
				AuthenticatorResponse: protocol.AuthenticatorResponse{
					ClientDataJSON: protocol.URLEncodedBase64(authenticatorClientDataJSON),
				},
				AuthenticatorData: protocol.URLEncodedBase64(authenticatorData),
				Signature:         protocol.URLEncodedBase64(authenticatorSignature),
				UserHandle:        protocol.URLEncodedBase64(authenticatorUserHandle),
			},
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

		fmt.Println("ParsedAssertionData: ", utils.PrintJSON(pca))

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
				Status:       "ok",
				ErrorMessage: "",
			},
		)
	}
}
