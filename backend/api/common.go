package api

type CommonResponse struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"errorMessage"`
}

type PublicKeyCredential struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}
