// Package router builds the top-level http.Handler for the forum.
package router

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	staticDir = "web/static"
	pagesDir  = "web/pages"
)

// New builds the mux. Nothing dynamic yet — Phase 1 is static files only.
func New() http.Handler {
	mux := http.NewServeMux()

	// Static assets: /static/css/style.css -> web/static/css/style.css
	fileServer := http.FileServer(http.Dir(staticDir))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Everything else is a "page" request, handled below so we control the
	// 404 response instead of leaking net/http's default file-not-found page.
	mux.HandleFunc("/", pageHandler)

	return mux
}

// pageHandler serves HTML pages out of web/pages. "/" maps to index.html;
// any other path maps to the matching file ("/login" or "/login.html" both
// resolve to web/pages/login.html). Anything not found gets a real custom
// 404 page with an actual 404 status code, not net/http's default text.
func pageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := strings.TrimPrefix(r.URL.Path, "/")
	if name == "" {
		name = "index"
	}
	name = strings.TrimSuffix(name, ".html")

	// filepath.Clean stops "../../etc/passwd" style paths escaping pagesDir.
	safeName := filepath.Clean("/" + name)
	fullPath := filepath.Join(pagesDir, safeName+".html")

	absPages, err := filepath.Abs(pagesDir)
	if err == nil {
		if absFull, err := filepath.Abs(fullPath); err != nil || !strings.HasPrefix(absFull, absPages) {
			serveNotFound(w, r)
			return
		}
	}

	if info, err := os.Stat(fullPath); err != nil || info.IsDir() {
		serveNotFound(w, r)
		return
	}

	http.ServeFile(w, r, fullPath)
}

// serveNotFound writes a real 404 status with our own page, instead of
// letting a handler fall through to net/http's plain-text default.
//
// We read the file and write it manually rather than calling http.ServeFile
// here, because ServeFile always sends its own 200 OK header on a
// successful open — calling it after we've already written 404 would either
// panic or silently get ignored as a "superfluous WriteHeader" call,
// corrupting the status our logging middleware records.
func serveNotFound(w http.ResponseWriter, r *http.Request) {
	notFoundPath := filepath.Join(pagesDir, "404.html")
	body, err := os.ReadFile(notFoundPath)
	if err != nil {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	w.Write(body)
}