package handler

type HTTPError struct {
	StatusCode int
	Message    string
}

func NewHTTPError(code int, msg string) *HTTPError {
	return &HTTPError{
		StatusCode: code,
		Message:    msg,
	}
}

func (e *HTTPError) Error() string {
	return e.Message
}
