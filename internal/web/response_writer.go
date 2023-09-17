package web

import "net/http"

// StubResponseWriter wraps the standard http.ResponseWriter allowing for more
// verbose logging
type StubResponseWriter struct {
	status int
	size   int
}

// Status provides an easy way to retrieve the status code
func (w *StubResponseWriter) Status() int {
	return w.status
}

// Size provides an easy way to retrieve the response size in bytes
func (w *StubResponseWriter) Size() int {
	return w.size
}

// Header returns & satisfies the http.ResponseWriter interface
func (w *StubResponseWriter) Header() http.Header {
	return make(http.Header)
}

// Write satisfies the http.ResponseWriter interface and
// captures data written, in bytes
func (w *StubResponseWriter) Write(_ []byte) (int, error) {
	return 0, nil
}

// WriteHeader satisfies the http.ResponseWriter interface and
// allows us to catch the status code
func (w *StubResponseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
}
