package payload

import "net/http"

type AppResponse struct {
	StatusCode int         `json:"code"`
	Data       interface{} `json:"data"`
	Message    string      `json:"message"`
}

// SuccessResponse simple success response
func SuccessResponse(data interface{}, message string) *AppResponse {
	return &AppResponse{http.StatusOK, data, message}
}