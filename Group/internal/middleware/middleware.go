// Package middleware holds HTTP middleware shared across the whole app.
package middleware

import (
	"log"
	"net/http"
	"time"
)

// statusRecorder wraps http.ResponseWriter so we can capture the status
// code a handler wrote (the standard library doesn't expose this).
type statusRecorder struct {
	http.ResponseWriter
	status  int
	written bool
}

func (r *statusRecorder) WriteHeader(code int) {
	// A handler may legitimately call WriteHeader more than once by mistake
	// (net/http logs it as "superfluous" and ignores it) — only the first
	// call is the real status, so ignore later ones instead of overwriting.
	if !r.written {
		r.status = code
		r.written = true
	}
	r.ResponseWriter.WriteHeader(code)
}

// Logging logs method, path, status code, and duration for every request.
// Wrap the top-level mux with this once in main.go.
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Default to 200: if the handler never calls WriteHeader explicitly
		// (e.g. it just calls w.Write), that's what net/http assumes too.
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)

		log.Printf("%s %s %d %s", r.Method, r.URL.Path, rec.status, time.Since(start))
	})
}