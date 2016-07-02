package server

import(
  "net/http"
)

type Handler interface {
  Render(path string, req *http.Request, res http.ResponseWriter, context map[string]interface{}) (bool, error)
  Passthrough(path string, req *http.Request, res http.ResponseWriter, context map[string]interface{}) (bool, error)
  Add(path string, handler Handler) error
  Walk(before string)
}