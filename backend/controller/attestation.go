package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
	goWebAuthn "github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"

	"leoho.io/it16th-webauthn-rp-server/api"
	"leoho.io/it16th-webauthn-rp-server/database"
	"leoho.io/it16th-webauthn-rp-server/utils"
	"leoho.io/it16th-webauthn-rp-server/webauthn"
)

var attestationSessionData *goWebAuthn.SessionData

func StartAttestationHandler(ctx *gin.Context) {
	fmt.Println("call /attestation/options")

	var request *api.CredentialCreationOptionsRequest
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

	user := &database.User{
		ID:          uuid.New().String(),
		Name:        request.Username,
		DisplayName: request.DisplayName,
	}
	excludeCredentialsOption := goWebAuthn.WithExclusions(user.CredentialExcludeList())
	options, sessionData, err := webauthn.WebAuthn.BeginRegistration(user, excludeCredentialsOption)
	if err != nil {
		fmt.Println("begin registration failed, error: ", err.Error())
		ctx.JSON(
			http.StatusInternalServerError,
			api.CommonResponse{
				Status:       "failed",
				ErrorMessage: "failed to create credential creation options, error: " + err.Error(),
			},
		)
		return
	}
	fmt.Println("begin registration success")

	user.Challenge = options.Response.Challenge.String()
	user.Credential = "`" + "{}" + "`"

	if err := database.CreateUser(user); err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			api.CommonResponse{
				Status:       "failed",
				ErrorMessage: "failed to create user, error: " + err.Error(),
			},
		)
		return
	}
	fmt.Println("create user success")

	attestationSessionData = sessionData

	ctx.JSON(
		http.StatusOK,
		api.CredentialCreationOptionsResponse{
			CommonResponse: api.CommonResponse{
				Status:       "success",
				ErrorMessage: "",
			},
			PublicKeyCredentialCreationOptions: options.Response,
		},
	)
}

func FinishAttestationHandler(ctx *gin.Context) {
	fmt.Println("call /attestation/result")

	var request *api.AuthenticatorAttestationResponseRequest
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
		fmt.Println("Get user by challenge success")

		ccr := protocol.CredentialCreationResponse{
			PublicKeyCredential: protocol.PublicKeyCredential{
				Credential: protocol.Credential{
					ID:   request.Id,
					Type: request.Type,
				},
				RawID:                  []byte(request.Id),
				ClientExtensionResults: request.GetClientExtensionResults,
			},
			AttestationResponse: protocol.AuthenticatorAttestationResponse{
				AttestationObject: protocol.URLEncodedBase64(request.Response.AttestationObject),
				AuthenticatorResponse: protocol.AuthenticatorResponse{
					ClientDataJSON: authenticatorClientDataJSON,
				},
			},
		}

		pcc, err := ccr.Parse()
		if err != nil {
			ctx.JSON(
				http.StatusBadRequest,
				api.CommonResponse{
					Status:       "failed",
					ErrorMessage: "failed to parse credential creation response, error: " + err.Error(),
				},
			)
			return
		}
		fmt.Println("Parse credential creation response success")

		credential, err := webauthn.WebAuthn.CreateCredential(foundUser, *attestationSessionData, pcc)
		if err != nil {
			ctx.JSON(
				http.StatusInternalServerError,
				api.CommonResponse{
					Status:       "failed",
					ErrorMessage: "failed to create credential, error: " + err.Error(),
				},
			)
			return
		}
		fmt.Println("Create credential success")

		credentialJSON, err := json.Marshal(credential)
		if err != nil {
			ctx.JSON(
				http.StatusInternalServerError,
				api.CommonResponse{
					Status:       "failed",
					ErrorMessage: "failed to marshal credential, error: " + err.Error(),
				},
			)
			return
		}
		fmt.Println("Marshal credential success")

		if err = database.UpdateUser(
			foundUser, database.User{
				Credential: "`" + string(credentialJSON) + "`",
			},
		); err != nil {
			ctx.JSON(
				http.StatusInternalServerError,
				api.CommonResponse{
					Status:       "failed",
					ErrorMessage: "failed to update user, error: " + err.Error(),
				},
			)
			return
		}
		fmt.Println("Update user credential success")

		ctx.JSON(
			http.StatusOK,
			api.CommonResponse{
				Status:       "success",
				ErrorMessage: "",
			},
		)
	}
}
