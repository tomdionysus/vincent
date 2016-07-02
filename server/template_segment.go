package server

import(
  "github.com/aymerick/raymond"
  "net/http"
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
func (me *TemplateSegment) Render(path string, req *http.Request, res http.ResponseWriter, context map[string]interface{}) (bool, error) {
  path = strings.TrimLeft(path,"/")
  
  if len(path) == 0 {
    // This is the last segment
    out, err := me.Template.Exec(context)
    if err!=nil { return false, err }
    res.Write([]byte(out))
    return true, nil
  }

  return me.Passthrough(path, req, res, context)
}