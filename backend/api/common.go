package api

type CommonResponse struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"errorMessage"`
}
