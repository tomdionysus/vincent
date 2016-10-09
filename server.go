package vincent

import (
	"crypto/tls"
	"github.com/aymerick/raymond"
	"github.com/tomdionysus/vincent/log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Server is a HTTP server for use with vincent projects.
type Server struct {
	Log log.Logger

	Root            *RouteSegment
	DefaultDocument string
}

// Return a new Server with the specified logger
func New(logger log.Logger) (*Server, error) {
	inst := &Server{
		Log:             logger,
		DefaultDocument: "index.html",
	}
	inst.Root = NewRouteSegment(inst)
	return inst, nil
}

// Walk the supplied basePath directory and parse all files and templates into routes
// using the route prefix specified.
func (svr *Server) LoadTemplates(routePrefix, basePath string) error {

	wfn := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		svr.Log.Debug("Loading: %s", path)

		ext := filepath.Ext(path)
		switch ext {
		case ".hbs":
			route := routePrefix + strings.TrimSuffix(path[len(basePath)+1:], ".hbs")
			template, err := raymond.ParseFile(path)
			if err != nil {
				return err
			}
			svr.Root.Add(route, NewTemplateSegment(template))
		case ".raw":
			fallthrough
		default:
			route := routePrefix + strings.TrimSuffix(path[len(basePath)+1:], ".raw")
			fn, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			svr.Root.Add(route, NewFileSegment(fn))
		}

		return nil
	}

	return filepath.Walk(basePath, wfn)
}

// Start the HTTP server on the specified address and port, of format "<host>:<port>", e.g. "localhost:8080"
func (svr *Server) Start(addr string) {
	go func() {
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			return
		}
		limitListener := NewConnLimitListener(250, ln.(*net.TCPListener))

		server := &http.Server{Handler: svr}

		server.Serve(limitListener)
	}()
}

// Start the HTTP server on the specified address and port, of format "<host>:<port>", e.g. "localhost:8080"
func (svr *Server) StartTLS(addr, certFile, keyFile string) {
	go func() {
		// TCP Layer
		tcpLn, err := net.Listen("tcp", addr)
		if err != nil {
			svr.Log.Error("Cannot Listen on %s", addr)
			return
		}

		// Conn limiter
		clLn := NewConnLimitListener(250, tcpLn.(*net.TCPListener))

		// TLS Layer
		cer, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			svr.Log.Error("Cannot Load Cert, Key %s, %s", certFile, keyFile)
			return
		}
		config := &tls.Config{Certificates: []tls.Certificate{cer}}

		tlsLn := tls.NewListener(clLn, config)

		server := &http.Server{Handler: svr}
		server.Serve(tlsLn)
	}()
}

func (svr *Server) AddController(path string, controller Controller) {
	svr.Root.AddController(path, controller)
}

// Support the http.Handler ServeHTTP method. This is called once per request
func (svr *Server) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	path := r.URL.EscapedPath()

	w := NewBufferedResponseWriter()
	t := time.Now()

	context := NewContext(svr, w, r)

	defer func() {
		rec := recover()
		size := formatByteSize(w.Buffer.Len())
		w.FlushToResponseWriter(wr)

		elapsed := time.Now().Sub(t).Seconds() / 1000
		svr.Log.Info("[%s] %s %s [%d] (%s/%.2fms)", r.RemoteAddr, r.Method, path, w.StatusCode, size, elapsed)

		if rec != nil {
			svr.Log.Error("> PANIC: %s", rec)
		}
	}()

	ok, err := svr.Root.Render(path, context)
	if err != nil {
		svr.Log.Error("Error while processing [%s] %s %s", r.Method, r.RemoteAddr, path)
		w.StatusCode = 500
		return
	}

	if !ok {
		w.StatusCode = 404
		return
	}
	w.StatusCode = 200
	return
}
