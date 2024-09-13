package api

type CredentialCreationOptionsRequest struct {
	Username               string                 `json:"username"`
	DisplayName            string                 `json:"displayName"`
	AuthenticatorSelection map[string]interface{} `json:"authenticatorSelection"`
	Attestation            string                 `json:"attestation"`
}

type CredentialCreationOptionsResponse struct {
	CommonResponse
	RP                     RelyingParty           `json:"rp"`
	User                   UserEntity             `json:"user"`
	Challenge              string                 `json:"challenge"`
	PubKeyCredParams       []PubKeyCredParam      `json:"pubKeyCredParams"`
	Timeout                int                    `json:"timeout"`
	ExcludeCredentials     []PublicKeyCredential  `json:"excludeCredentials"`
	AuthenticatorSelection map[string]interface{} `json:"authenticatorSelection"`
	Attestation            string                 `json:"attestation"`
	Extensions             map[string]interface{} `json:"extensions"`
}

type RelyingParty struct {
	Name string `json:"name"`
}

type UserEntity struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

type PubKeyCredParam struct {
	Type string `json:"type"`
	Alg  int    `json:"alg"`
}

type AuthenticatorAttestationResponseRequest struct {
	Id                        string                           `json:"id"`
	Response                  AuthenticatorAttestationResponse `json:"response"`
	GetClientExtensionResults map[string]interface{}           `json:"getClientExtensionResults"`
	Type                      string                           `json:"type"`
}

type AuthenticatorAttestationResponse struct {
	AttestationObject string `json:"attestationObject"`
	ClientDataJSON    string `json:"clientDataJSON"`
}
