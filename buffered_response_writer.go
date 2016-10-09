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
func (brw *BufferedResponseWriter) Header() http.Header {
	return brw.Headers
}

// Write the supplied data
func (brw *BufferedResponseWriter) Write(data []byte) (int, error) {
	return brw.Buffer.Write(data)
}

// Set the HTTP status code
func (brw *BufferedResponseWriter) WriteHeader(code int) {
	brw.StatusCode = code
}

// Write HTTP Status, Headers and Flush all data to the supplied http.ResponseWriter
func (brw *BufferedResponseWriter) FlushToResponseWriter(w http.ResponseWriter) error {
	outheader := w.Header()
	for k, v := range brw.Headers {
		outheader[k] = v
	}
	w.WriteHeader(brw.StatusCode)
	_, err := brw.Buffer.WriteTo(w)
	return err
}
