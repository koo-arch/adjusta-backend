package errors

import "net/http"

type APIError struct {
	StatusCode int 			  `json:"-"`	// "-" means that this field will not be marshaled
	Message string 			  `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

func NewAPIError(statusCode int, Message string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message: Message,
	}
}

func NewValidationError(details map[string]string) *APIError {
	return &APIError{
		StatusCode: http.StatusBadRequest,
		Message: "送信に失敗しました",
		Details: details,
	}
}

func (e *APIError) Error() string {
	return e.Message
}

var (
	InternalErrorMessage = "サーバーでエラーが発生しました"
)