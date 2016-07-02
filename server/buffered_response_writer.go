package server

import(
  "net/http"
  "bytes"
)

type BufferedResponseWriter struct {
  Headers http.Header
  Buffer *bytes.Buffer
  StatusCode int
}

func NewBufferedResponseWriter() *BufferedResponseWriter {
  return &BufferedResponseWriter{
    Buffer: &bytes.Buffer{},
    Headers: http.Header{},
  }
}

func (me *BufferedResponseWriter) Header() http.Header {
  return me.Headers
}

func (me *BufferedResponseWriter) Write(data []byte) (int, error) {
  return me.Buffer.Write(data)
}

func (me *BufferedResponseWriter) WriteHeader(code int) {
  me.StatusCode = code
}

func (me *BufferedResponseWriter) WriteToResponseWriter(w http.ResponseWriter) error {
  outheader := w.Header()
  for k,v := range me.Headers { outheader[k] = v }
  w.WriteHeader(me.StatusCode)
  _, err := me.Buffer.WriteTo(w)
  return err
}
