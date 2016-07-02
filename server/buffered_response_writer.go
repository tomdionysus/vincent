package server

import(
  "net/http"
  "bytes"
  "encoding/gob"
)

type BufferedResponseWriter struct {
  Buffer *bytes.Buffer
  Headers http.Header
  StatusCode int
}

func NewBufferedResponseWriter() *BufferedResponseWriter {
  return &BufferedResponseWriter{
    Buffer: &bytes.Buffer{},
    Headers: http.Header{},
  }
}

func DeserialiseBufferedResponseWriter(buf []byte) *BufferedResponseWriter {
  me := &BufferedResponseWriter{}

  dec := gob.NewDecoder(bytes.NewBuffer(buf))
  dec.Decode(&me.StatusCode)
  dec.Decode(&me.Headers)

  b := []byte{}
  dec.Decode(&b)
  me.Buffer = bytes.NewBuffer(b)

  return me
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

func (me *BufferedResponseWriter) Serialise() []byte {
  buf := &bytes.Buffer{}
  enc := gob.NewEncoder(buf)
  enc.Encode(me.StatusCode)
  enc.Encode(me.Headers)
  enc.Encode(me.Buffer.Bytes())
  return buf.Bytes()
}

func (me *BufferedResponseWriter) WriteToResponseWriter(w http.ResponseWriter) error {
  outheader := w.Header()
  for k,v := range me.Headers { outheader[k] = v }
  w.WriteHeader(me.StatusCode)
  _, err := me.Buffer.WriteTo(w)
  return err
}
