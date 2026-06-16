package main

import (
	"fmt"
	"net/http"

	"groupie-tracker/internal/handlers"
)

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/api/search", handlers.SearchAPI)
	http.HandleFunc("/api/filter", handlers.FilterAPI)
	http.HandleFunc("/artist/", handlers.ArtistPage)
	http.HandleFunc("/notfound", handlers.NotFoundPage)
	http.HandleFunc("/", handlers.Home)

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}