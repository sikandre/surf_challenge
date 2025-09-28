package error

type APIError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e *APIError) Error() string {
	return e.Message
}

func NewAPIError(message string, code int) *APIError {
	return &APIError{
		Message: message,
		Code:    code,
	}
}
