package server

import(
  "github.com/aymerick/raymond"
  "github.com/tomdionysus/vincent"
  "github.com/tomdionysus/vincent/log"
  "path/filepath"
  "os"
  "fmt"
  "net/http"
  "time"
  "strings"
)

type Server struct {
  Log log.Logger

  Root *RouteSegment
}

func New(logger log.Logger) (*Server, error) {
  inst := &Server{
    Log: log.NewScopedLogger("Vincent", logger),
  }
  inst.Root = NewRouteSegment(inst)
  return inst, nil
}

func (me *Server) LoadTemplates(routePrefix, basePath string) error {

  wfn := func(path string, info os.FileInfo, err error) error {
    if info.IsDir() { return nil }
    me.Log.Debug("Loading: %s", path)

    ext := filepath.Ext(path)
    switch ext {
    case ".hbs":
      route := routePrefix+strings.TrimSuffix(path[len(basePath)+1:], ".hbs")
      template, err := raymond.ParseFile(path)
      if err!=nil { return err }
      me.Root.Add(route, NewTemplateSegment(template))
    default:
      route := routePrefix+strings.TrimSuffix(path[len(basePath)+1:], ".hbs")
      fn, err := filepath.Abs(path)
      if err!=nil { return err }
      me.Root.Add(route, NewFileSegment(fn))
    }

    return nil
  }

 return filepath.Walk(basePath, wfn)
}

func (me *Server) Start(port uint16) {
  go func(){ 
    http.ListenAndServe(fmt.Sprintf(":%d",port), me) 
  }()
}

func (me *Server) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
  path := r.URL.EscapedPath()

  w := NewBufferedResponseWriter()

  t := time.Now()

  context := map[string]interface{}{
    "vincent": map[string]interface{}{
      "version": vincent.VERSION,
      "port": 8080,
    },
  }

  defer func(){
    rec := recover();
    w.WriteToResponseWriter(wr)

    elapsed := time.Now().Sub(t).Seconds() / 1000
    size := formatByteSize(w.Buffer.Len())
    me.Log.Info("[%s] %s %s [%d] (%s/%.2fms)", r.RemoteAddr, r.Method, path, w.StatusCode, size, elapsed)

    if rec != nil { me.Log.Error("> PANIC: %s", rec) }
  }()

  ok, err := me.Root.Render(path, r, w, context)
  if err!=nil {
    me.Log.Error("Error while processing [%s] %s %s", r.Method, r.RemoteAddr, path)
    w.StatusCode = 500
    return
  }

  if !ok { w.StatusCode = 404; return }
  w.StatusCode = 200
  return

}