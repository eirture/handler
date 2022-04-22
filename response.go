package shandler

import (
	"io"
)

type Response struct {
	StatusCode int
	Body       io.ReadCloser
}

func (r *Response) WithStatusCode(code int) *Response {
	r.StatusCode = code
	return r
}

func (r *Response) WithBodyReader(reader io.ReadCloser) *Response {
	r.Body = reader
	return r
}
