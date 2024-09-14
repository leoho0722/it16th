package api

import "github.com/go-webauthn/webauthn/protocol"

type CredentialCreationOptionsRequest struct {
	Username               string                 `json:"username"`
	DisplayName            string                 `json:"displayName"`
	AuthenticatorSelection map[string]interface{} `json:"authenticatorSelection"`
	Attestation            string                 `json:"attestation"`
}

type CredentialCreationOptionsResponse struct {
	CommonResponse
	protocol.PublicKeyCredentialCreationOptions
}

type AuthenticatorAttestationResponseRequest struct {
	// Id CredentialID from Authneticator
	Id                        string                           `json:"id"`
	Response                  AuthenticatorAttestationResponse `json:"response"`
	GetClientExtensionResults map[string]interface{}           `json:"getClientExtensionResults"`
	Type                      string                           `json:"type"`
}

type AuthenticatorAttestationResponse struct {
	AttestationObject string `json:"attestationObject"`
	ClientDataJSON    string `json:"clientDataJSON"`
}
