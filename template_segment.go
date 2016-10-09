package vincent

import (
	"github.com/aymerick/raymond"
	"strings"
)

// A Segment representing a handlebars template
type TemplateSegment struct {
	RouteSegment
	Template *raymond.Template
}

// Return a new TemplateSegment with the supplied raymond.Template
func NewTemplateSegment(template *raymond.Template) *TemplateSegment {
	inst := &TemplateSegment{
		Template: template,
	}
	return inst
}

// If the path ends with this segment, render the template using the supplied context to the responsewriter.
// Otherwise, passthrough to sub segments.
func (tsg *TemplateSegment) Render(path string, context *Context) (bool, error) {
	ok, err := tsg.CallControllers(context)
	if !ok || err != nil {
		return ok, err
	}

	path = strings.TrimLeft(path, "/")

	if len(path) == 0 {
		// This is the last segment
		out, err := tsg.Template.Exec(context.Output)
		if err != nil {
			return false, err
		}
		context.ResponseWriter.Write([]byte(out))
		return true, nil
	}

	return tsg.Passthrough(path, context)
}
