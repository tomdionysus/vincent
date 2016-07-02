package server

import(
  "net/http"
  "strings"
  "io/ioutil"
  "mime"
  "path/filepath"
)

// A segment of a route that represents a raw file
type FileSegment struct {
  RouteSegment
  Filename string
}

// Return a new FileSegment with the supplied filename
func NewFileSegment(filename string) *FileSegment {
  inst := &FileSegment{
    Filename: filename,
  }
  return inst
}

// If the path refers to this segment, render the supplied path to the responsewriter. Otherwise, passthrough to
// sub segments.
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