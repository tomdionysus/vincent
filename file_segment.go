package vincent

import (
	"io/ioutil"
	"mime"
	"path/filepath"
	"strings"
)

// A segment of a route that repcontext.ResponseWriterents a raw file
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

// If the path refers to this segment, render the supplied path to the context.ResponseWriterponsewriter. Otherwise, passthrough to
// sub segments.
func (fsg *FileSegment) Render(path string, context *Context) (bool, error) {
	ok, err := fsg.CallControllers(context)
	if !ok || err != nil {
		return ok, err
	}

	path = strings.TrimLeft(path, "/")

	if len(path) == 0 {
		// This is the last segment

		ext := filepath.Ext(fsg.Filename)
		h := context.ResponseWriter.Header()

		h["Content-Type"] = append(h["Content-Type"], mime.TypeByExtension(ext))

		out, err := ioutil.ReadFile(fsg.Filename)
		if err != nil {
			return false, err
		}
		context.ResponseWriter.Write(out)
		return true, nil
	}

	return fsg.Passthrough(path, context)
}
