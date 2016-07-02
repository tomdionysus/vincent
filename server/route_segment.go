package server

import(
  "net/http"
  "strings"
)

// Represents a single segment of a route, e.g. http://host:port/<segment>/<segment>...
type RouteSegment struct {
  Server *Server
  Segments map[string]Handler
}

// Return a new RouteSegment for the supplied server
func NewRouteSegment(svr *Server) *RouteSegment {
  inst := &RouteSegment{
    Server: svr,
    Segments: map[string]Handler{},
  }
  return inst
}

// Parse the path and passthrough, substituting the DefaultDocument (if available) for an empty segment.
func (me *RouteSegment) Render(path string, req *http.Request, res http.ResponseWriter, context map[string]interface{}) (bool, error) {
  path = strings.TrimLeft(path,"/")
  // Special case if path is empty.
  if sgm, ok := me.Segments[me.Server.DefaultDocument]; path == "" && ok { return sgm.Render(path, req, res, context) }
  if len(path)!=0 { return me.Passthrough(path, req, res, context) }
  return false, nil
}

// Process the path and call Render on subroute handlers
func (me *RouteSegment) Passthrough(path string, req *http.Request, res http.ResponseWriter, context map[string]interface{}) (bool, error) {
  c := strings.Index(path,"/")

  var sgmName string
  if c == -1 {
    sgmName = path
    path = ""
  } else {
    sgmName = path[:c]
    path = path[c+1:]
  }

  sgm, ok := me.Segments[sgmName]
  if ok {
    return sgm.Render(path, req, res, context)
  } else {
    
    // Otherwise not found
    return false, nil
  }
}

// Add a route and handler to this segment at this segment's path
func (me *RouteSegment) Add(path string, handler Handler) error {
  path = strings.TrimLeft(path,"/")
  c := strings.Index(path,"/")

  var sgmName string
  if c == -1 {
    sgmName = path
    path = ""
  } else {
    sgmName = path[:c]
    path = path[c+1:]
  }

  sgm, ok := me.Segments[sgmName]
  if !ok {
    if len(path) == 0 {
      me.Segments[sgmName] = handler
    } else {
      sgm = NewRouteSegment(me.Server)
      me.Segments[sgmName] = sgm
      sgm.Add(path, handler)
    }
  } else {
    sgm.Add(path, handler)
  }

  return nil
}

// A function to walk the segment tree
type RouteSegmentWalkFunc func(path, name string) bool

// Walk the segment tree calling the supplied RouteSegmentWalkFunc for each possible route leaf
func (me *RouteSegment) Walk(path string, fn RouteSegmentWalkFunc) bool {
  for name, segment := range me.Segments {
    if !fn(path, name) || segment.Walk(path+"/"+name, fn) { return false }
  }
  return true
}