package vincent

import (
	"bytes"
	"net/http"
)

// A ResponseWriter that buffers everything written to it
type BufferedResponseWriter struct {
	Headers    http.Header
	Buffer     *bytes.Buffer
	StatusCode int
}

// Return a new BufferedResponseWriter
func NewBufferedResponseWriter() *BufferedResponseWriter {
	return &BufferedResponseWriter{
		Buffer:  &bytes.Buffer{},
		Headers: http.Header{},
	}
}

// Return the current headers
func (me *BufferedResponseWriter) Header() http.Header {
	return me.Headers
}

// Write the supplied data
func (me *BufferedResponseWriter) Write(data []byte) (int, error) {
	return me.Buffer.Write(data)
}

// Set the HTTP status code
func (me *BufferedResponseWriter) WriteHeader(code int) {
	me.StatusCode = code
}

// Write HTTP Status, Headers and Flush all data to the supplied http.ResponseWriter
func (me *BufferedResponseWriter) FlushToResponseWriter(w http.ResponseWriter) error {
	outheader := w.Header()
	for k, v := range me.Headers {
		outheader[k] = v
	}
	w.WriteHeader(me.StatusCode)
	_, err := me.Buffer.WriteTo(w)
	return err
}
