package api

import "github.com/go-webauthn/webauthn/protocol"

type CredentialGetOptionsRequest struct {
	Username         string `json:"username"`
	UserVerification string `json:"userVerification"`
}

type CredentialGetOptionsResponse struct {
	CommonResponse
	protocol.PublicKeyCredentialRequestOptions
}

type AuthenticatorAssertionResponseRequest struct {
	Id                        string                         `json:"id"`
	Response                  AuthenticatorAssertionResponse `json:"response"`
	GetClientExtensionResults map[string]interface{}         `json:"getClientExtensionResults"`
	Type                      string                         `json:"type"`
}

type AuthenticatorAssertionResponse struct {
	ClientDataJSON    string `json:"clientDataJSON"`
	AuthenticatorData string `json:"authenticatorData"`
	Signature         string `json:"signature"`
	UserHandle        string `json:"userHandle,omitempty"`
}
