package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"groupie-tracker/internal/api"
)

type PageData struct {
	Artist    api.Artist
	Relations map[string][]string
}

type ErrorData struct {
	Code    int
	Title   string
	Message string
}

func renderError(w http.ResponseWriter, code int, title, message string) {
	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		http.Error(w, message, code)
		return
	}
	w.WriteHeader(code)
	tmpl.Execute(w, ErrorData{
		Code:    code,
		Title:   title,
		Message: message,
	})
}

func home(w http.ResponseWriter, r *http.Request) {
	artists, err := api.GetArtists()
	if err != nil {
		renderError(w, 500, "Something went wrong", "We couldn't load the artists. Please try again.")
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		renderError(w, 500, "Something went wrong", "Failed to load the page template.")
		return
	}

	tmpl.Execute(w, artists)
}

func artistPage(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/artist/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		renderError(w, 404, "Artist not found", "The artist you're looking for doesn't exist.")
		return
	}

	artists, err := api.GetArtists()
	if err != nil {
		renderError(w, 500, "Something went wrong", "We couldn't load the artists. Please try again.")
		return
	}

	var found *api.Artist
	for _, a := range artists {
		if a.ID == id {
			found = &a
			break
		}
	}

	if found == nil {
		renderError(w, 404, "Artist not found", "No artist with that ID exists.")
		return
	}

	relations, err := api.GetRelations()
	if err != nil {
		renderError(w, 500, "Something went wrong", "We couldn't load the tour dates.")
		return
	}

	tmpl, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		renderError(w, 500, "Something went wrong", "Failed to load the page template.")
		return
	}

	tmpl.Execute(w, PageData{
		Artist:    *found,
		Relations: relations[id].DatesLocations,
	})
}

func searchAPI(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	artists, err := api.GetArtists()
	if err != nil {
		http.Error(w, "Failed to load artists", http.StatusInternalServerError)
		return
	}

	results := api.SearchArtists(query, artists)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func filterAPI(w http.ResponseWriter, r *http.Request) {
	minYear, _ := strconv.Atoi(r.URL.Query().Get("minYear"))
	maxYear, _ := strconv.Atoi(r.URL.Query().Get("maxYear"))
	members, _ := strconv.Atoi(r.URL.Query().Get("members"))

	artists, err := api.GetArtists()
	if err != nil {
		http.Error(w, "Failed to load artists", http.StatusInternalServerError)
		return
	}

	results := api.FilterArtists(artists, minYear, maxYear, members)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/api/search", searchAPI)
	http.HandleFunc("/api/filter", filterAPI)
	http.HandleFunc("/artist/", artistPage)
	http.HandleFunc("/", home)

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}