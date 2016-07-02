package server

import(
  "github.com/aymerick/raymond"
  "github.com/tomdionysus/vincent"
  "github.com/tomdionysus/vincent/log"
  "path/filepath"
  "os"
  "time"
  "strings"
  "net"
  "net/http"
)

// Server is a HTTP server for use with vincent projects. 
type Server struct {
  Log log.Logger

  Root *RouteSegment
  DefaultDocument string
}

// Return a new Server with the specified logger
func New(logger log.Logger) (*Server, error) {
  inst := &Server{
    Log: logger,
    DefaultDocument: "index.html",
  }
  inst.Root = NewRouteSegment(inst)
  return inst, nil
}

// Walk the supplied basePath directory and parse all files and templates into routes
// using the route prefix specified.
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
    case ".raw":
      fallthrough
    default:
      route := routePrefix+strings.TrimSuffix(path[len(basePath)+1:], ".raw")
      fn, err := filepath.Abs(path)
      if err!=nil { return err }
      me.Root.Add(route, NewFileSegment(fn))
    }

    return nil
  }

 return filepath.Walk(basePath, wfn)
}

// Start the HTTP server on the specified address and port, of format "<host>:<port>", e.g. "localhost:8080"
func (me *Server) Start(addr string) {
  go func(){ 
    ln, err := net.Listen("tcp", addr)
    if err != nil { return }
    limitListener := NewConnLimitListener(250, ln.(*net.TCPListener))

    server := &http.Server{ Handler: me }

    server.Serve(limitListener)
  }()
}

// Support the http.Handler ServeHTTP method. This is called once per request
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
    size := formatByteSize(w.Buffer.Len())
    w.FlushToResponseWriter(wr)

    elapsed := time.Now().Sub(t).Seconds() / 1000
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