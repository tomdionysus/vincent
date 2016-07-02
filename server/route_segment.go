package server

import(
  "net/http"
  "strings"
  "fmt"
)

type RouteSegment struct {
  Server *Server
  Segments map[string]Handler
}

func NewRouteSegment(svr *Server) *RouteSegment {
  inst := &RouteSegment{
    Server: svr,
    Segments: map[string]Handler{},
  }
  return inst
}

func (me *RouteSegment) Render(path string, req *http.Request, res http.ResponseWriter, context map[string]interface{}) (bool, error) {
  path = strings.TrimLeft(path,"/")
  if len(path)!=0 { return me.Passthrough(path, req, res, context) }
  return false, nil
}
  
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
    return false, nil
  }
}

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

func (me *RouteSegment) Walk(before string) {
  for k, v := range me.Segments {
    pth := fmt.Sprintf("Tree: %s/%s",before,k)
    me.Server.Log.Info(pth)
    v.Walk(pth)
  }
}