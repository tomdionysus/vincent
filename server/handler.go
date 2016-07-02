package server

import(
  "net/http"
)

// An interface that handles a Route segment
type Handler interface {
  Render(path string, req *http.Request, res http.ResponseWriter, context map[string]interface{}) (bool, error)
  Passthrough(path string, req *http.Request, res http.ResponseWriter, context map[string]interface{}) (bool, error)
  Add(path string, handler Handler) error
  Walk(path string, fn RouteSegmentWalkFunc) bool
}

// A function that handle a route segment
type HandlerFunc func(path string, req *http.Request, res http.ResponseWriter, context map[string]interface{}) (bool, error)