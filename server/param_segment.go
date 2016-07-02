package server

import(
  "net/http"
  "strings"
)

// A segment of a route that represents a parameter
type ParamSegment struct {
  RouteSegment
  ParamName string
}

// Return a new ParamSegment with the supplied name, e.g. "identity" or "identity.name"
func NewParamSegment(paramName string) *ParamSegment {
  inst := &ParamSegment{
    ParamName: paramName,
  }
  return inst
}

// Load the current segment value into the context and passthrough.
func (me *ParamSegment) Render(path string, req *http.Request, res http.ResponseWriter, context map[string]interface{}) (bool, error) {
  path = strings.TrimLeft(path,"/")
  
  c := strings.Index(path,"/")

  var value string
  if c == -1 {
    value = path
    path = ""
  } else {
    value = path[:c]
    path = path[c+1:]
  }

  context[me.ParamName] = value

  return me.Passthrough(path, req, res, context)
}