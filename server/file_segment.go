package server

import(
  "net/http"
  "strings"
  "io/ioutil"
  "mime"
  "path/filepath"
)

type FileSegment struct {
  RouteSegment
  Filename string
}

func NewFileSegment(filename string) *FileSegment {
  inst := &FileSegment{
    Filename: filename,
  }
  return inst
}

func (me *FileSegment) Render(path string, req *http.Request, res http.ResponseWriter, context map[string]interface{}) (bool, error) {
  path = strings.TrimLeft(path,"/")
  
  if len(path) == 0 {
    // This is the last segment

    ext := filepath.Ext(me.Filename)
    h := res.Header()

    h["Content-Type"] = append(h["Content-Type"], mime.TypeByExtension(ext))

    out, err := ioutil.ReadFile(me.Filename)
    if err!=nil { return false, err }
    res.Write(out)
    return true, nil
  }

  return me.Passthrough(path, req, res, context)
}