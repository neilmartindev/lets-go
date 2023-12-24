package main

import (
	"html/template"
	"path/filepath"

	"snippetbox.martinlabs.io/internal/models"
)

// Define a templateData type to act as the holding
// structure for any dynamic data that we want to pass to our HTML
// templates.

type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
