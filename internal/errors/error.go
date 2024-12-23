package errors

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

func (e *APIError) Error() string {
	return e.Message
}