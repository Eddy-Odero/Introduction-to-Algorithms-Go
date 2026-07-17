// Package router builds the top-level http.Handler for the forum.
package router

import (
	"net/http"

	"forum/internal/handlers"
)

const staticDir = "web/static"

// New builds the mux. Every page is rendered by a handler in
// internal/handlers — no generic static-file-to-page mapping anymore,
// since pages are now Go templates with real data, not flat HTML files.
func New() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(staticDir))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			handlers.NotFound(w, r)
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handlers.Home(w, r)
	})

	get(mux, "/login", handlers.Login)
	get(mux, "/register", handlers.Register)
	get(mux, "/new-post", handlers.NewPost)
	get(mux, "/post", handlers.PostDetail)
	get(mux, "/category", handlers.Category)
	get(mux, "/profile", handlers.Profile)

	return mux
}

// get registers a GET-only route, rejecting other methods with 405 instead
// of silently falling through to whatever the mux would otherwise do.
func get(mux *http.ServeMux, pattern string, h http.HandlerFunc) {
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h(w, r)
	})
}