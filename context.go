package vincent

import (
	"net/http"
)

type Context struct {
	Server *Server

	Request        *http.Request
	ResponseWriter http.ResponseWriter

	Params map[string]interface{}

	Input  map[string]interface{}
	Output map[string]interface{}
}

func NewContext(server *Server, w http.ResponseWriter, r *http.Request) *Context {
	i := &Context{
		Server:         server,
		Request:        r,
		ResponseWriter: w,
		Params:         map[string]interface{}{},
		Input:          map[string]interface{}{},
		Output:         map[string]interface{}{},
	}
	return i
}
