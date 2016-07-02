package server

import(
  "github.com/aymerick/raymond"
  "net/http"
  "strings"
)

type TemplateSegment struct {
  RouteSegment
  Template *raymond.Template
}

func NewTemplateSegment(template *raymond.Template) *TemplateSegment {
  inst := &TemplateSegment{
    Template: template,
  }
  return inst
}

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