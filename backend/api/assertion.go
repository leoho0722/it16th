package api

type CredentialGetOptions struct {
	Username         string `json:"username"`
	UserVerification string `json:"userVerification"`
}

type CredentialGetOptionsResponse struct {
	CommonResponse
	Challenge        string                `json:"challenge"`
	Timeout          int                   `json:"timeout"`
	RPId             string                `json:"rpId"`
	AllowCredentials []PublicKeyCredential `json:"allowCredentials"`
	UserVerification string                `json:"userVerification"`
}

type AuthenticatorAssertionResponseRequest struct {
	Id                        string                         `json:"id"`
	Response                  AuthenticatorAssertionResponse `json:"response"`
	GetClientExtensionResults map[string]interface{}         `json:"getClientExtensionResults"`
	Type                      string                         `json:"type"`
}

type AuthenticatorAssertionResponse struct {
	AuthenticatorData string `json:"authenticatorData"`
	ClientDataJSON    string `json:"clientDataJSON"`
	Signature         string `json:"signature"`
	UserHandle        string `json:"userHandle"`
}
