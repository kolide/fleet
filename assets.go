package main

import (
	"net/http"
	"strings"
	"html/template"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/contrib/renders/multitemplate"
)

// Helper functions to make Gin able to load templates from go-bindata-assetfs
type binaryFileSystem struct {
	fs http.FileSystem
}

func (b *binaryFileSystem) Open(name string) (http.File, error) {
	return b.fs.Open(name)
}

func (b *binaryFileSystem) Exists(prefix string, filepath string) bool {

	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		if _, err := b.fs.Open(p); err != nil {
			return false
		}
		return true
	}
	return false
}

func BinaryFileSystem(root string) *binaryFileSystem {
	fs := &assetfs.AssetFS{Asset, AssetDir, AssetInfo, root}
	return &binaryFileSystem{
		fs,
	}
}

//
//
func loadTemplates(list ...string) multitemplate.Render {
	r := multitemplate.New()
  for _, x := range list {
	  templateString, err := Asset("frontend/templates/" + x)
		if err != nil {
			panic(err)
		}
    tmplMessage, err := template.New(x).Parse(string(templateString))
		if err != nil {
      panic(err)
    }
    r.Add(x, tmplMessage)
  }
  return r
}
