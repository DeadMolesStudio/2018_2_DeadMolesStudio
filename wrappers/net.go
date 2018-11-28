package wrappers

import (
	"net/http"
)

type StatusWrapperForResponseWriter struct {
	http.ResponseWriter
	Status int
}

func (w *StatusWrapperForResponseWriter) WriteHeader(code int) {
	w.Status = code
	w.ResponseWriter.WriteHeader(code)
}
