// Package view embeds the HTML templates and static assets.
package view

import (
	"embed"
	"html/template"
	"io/fs"
)

//go:embed templates/*.html static/*
var files embed.FS

// Pages parses layout.html together with each page so that every page can
// define its own "content" and "hero" blocks.
func Pages(names ...string) (map[string]*template.Template, error) {
	out := map[string]*template.Template{}
	for _, n := range names {
		t, err := template.ParseFS(files, "templates/layout.html", "templates/"+n+".html")
		if err != nil {
			return nil, err
		}
		out[n] = t
	}
	return out, nil
}

func Static() fs.FS {
	sub, err := fs.Sub(files, "static")
	if err != nil {
		panic(err)
	}
	return sub
}
