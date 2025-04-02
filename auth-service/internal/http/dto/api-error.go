package dto

type ApiError struct {
	Message string `json:"message"`
}

func NewApiError(message string) *ApiError {
	return &ApiError{Message: message}
}
