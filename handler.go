package vincent
import(
  "net/http"
)

// An interface that handles a Route segment
type Handler interface {
  Render(path string, context *Context) (bool, error)
  Passthrough(path string, context *Context) (bool, error)
  Add(path string, handler Handler) error
  AddController(path string, controller Controller)
  CallControllers(context *Context) (bool, error)
  Walk(path string, fn RouteSegmentWalkFunc) bool
}

// A function that handle a route segment
type HandlerFunc func(path string, req *http.Request, res http.ResponseWriter, context map[string]interface{}) (bool, error)