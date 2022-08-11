package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

const (
	DefaultHTTPErrorStatusCode = http.StatusInternalServerError
)

var (
	DefaultHandler = NewHandler()
)

type HandlerFunc func(rw http.ResponseWriter, req *http.Request) interface{}

type ErrMarshalFn func(err error) (statusCode int, body []byte)

type respError struct {
	Error string `json:"error"`
}

func DefaultErrMarshalFn(err error) (sc int, b []byte) {
	sc = DefaultHTTPErrorStatusCode
	if hErr, ok := err.(*HTTPError); ok {
		sc = hErr.StatusCode
	}

	b, mErr := json.Marshal(&respError{Error: err.Error()})
	if mErr != nil {
		b = []byte(err.Error())
	}
	return
}

func DefaultWriteErrHandlerFn(n int64, err error) {}

type Handler struct {
	errMarshalFn      ErrMarshalFn
	writeErrHandlerFn func(n int64, err error)
}

func NewHandler() *Handler {
	return &Handler{
		errMarshalFn:      DefaultErrMarshalFn,
		writeErrHandlerFn: DefaultWriteErrHandlerFn,
	}
}

func (h *Handler) Wrap(hfn HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rst := hfn(rw, req)

		var (
			statusCode = http.StatusOK
			body       = io.NopCloser(bytes.NewReader([]byte{}))
		)

		switch v := rst.(type) {
		case error:
			var b []byte
			statusCode, b = h.errMarshalFn(v)
			body = io.NopCloser(bytes.NewReader(b))
		case *Response:
			statusCode = v.StatusCode
			body = v.Body
		default:
		}

		defer body.Close()
		rw.WriteHeader(statusCode)
		if n, err := io.Copy(rw, body); err != nil {
			h.writeErrHandlerFn(n, err)
		}
	}
}

func (h *Handler) SetErrMarshalFn(fn ErrMarshalFn) *Handler {
	if fn == nil {
		fn = DefaultErrMarshalFn
	}

	h.errMarshalFn = fn
	return h
}

func (h *Handler) SetWriteErrHandlerFn(fn func(int64, error)) *Handler {
	if fn == nil {
		fn = DefaultWriteErrHandlerFn
	}
	h.writeErrHandlerFn = fn
	return h
}

func WrapHandler(h HandlerFunc) http.HandlerFunc {
	return DefaultHandler.Wrap(h)
}
