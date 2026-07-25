// Package view loads and renders the site's html/template pages.
//
// Each page is parsed as base.html + <page>.html together, kept in its own
// *template.Template. We deliberately do NOT use template.ParseGlob on the
// whole directory: every page defines a block named "content", and with
// ParseGlob all files share one namespace — the last file parsed would
// silently win and every page would render identical content. Parsing each
// page as its own {base, page} pair avoids that entirely.
package view

import (
	"bytes"
	"html/template"
	"net/http"
	"path/filepath"
)

var pages = map[string]*template.Template{}

// pageNames lists every template in web/templates that pairs with base.html.
var pageNames = []string{
	"index", "login", "register", "new-post", "post", "category", "profile",
	"404", "500",
}

// Load parses every page template. Call once at startup; if it returns an
// error the server should fail fast rather than serve broken pages.
func Load(dir string) error {
	base := filepath.Join(dir, "base.html")
	partials := filepath.Join(dir, "partials.html")

	for _, name := range pageNames {
		t, err := template.New("base.html").ParseFiles(base, partials, filepath.Join(dir, name+".html"))
		if err != nil {
			return err
		}
		pages[name] = t
	}
	return nil
}

// Render executes the named page's "base" block with data into a buffer
// first. Rendering directly to w would mean a mid-template error leaves a
// half-written response with a status code we already committed to — by
// buffering, a failure can still become a clean 500 instead.
func Render(w http.ResponseWriter, status int, name string, data any) {
	t, ok := pages[name]
	if !ok {
		http.Error(w, "template not found: "+name, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "base", data); err != nil {
		http.Error(w, "internal error rendering page", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}