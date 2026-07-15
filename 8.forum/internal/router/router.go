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

func New() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(staticDir))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	mux.HandleFunc("/", pageHandler)

	return mux
}

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