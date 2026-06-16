package handlers

import (
	"encoding/json"
	"groupie-tracker/internal/api"
	"html/template"
	"net/http"
	"strconv"
	"strings"
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

func RenderError(w http.ResponseWriter, code int, title, message string) {
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

func Home(w http.ResponseWriter, r *http.Request) {
	artists, err := api.GetArtists()
	if err != nil {
		RenderError(w, 500, "Something went wrong", "We couldn't load the artists. Please try again.")
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		RenderError(w, 500, "Something went wrong", "Failed to load the page template.")
		return
	}

	tmpl.Execute(w, artists)
}

func ArtistPage(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/artist/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		RenderError(w, 404, "Artist not found", "The artist you're looking for doesn't exist.")
		return
	}

	artists, err := api.GetArtists()
	if err != nil {
		RenderError(w, 500, "Something went wrong", "We couldn't load the artists. Please try again.")
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
		RenderError(w, 404, "Artist not found", "No artist with that ID exists.")
		return
	}

	relations, err := api.GetRelations()
	if err != nil {
		RenderError(w, 500, "Something went wrong", "We couldn't load the tour dates.")
		return
	}

	tmpl, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		RenderError(w, 500, "Something went wrong", "Failed to load the page template.")
		return
	}

	tmpl.Execute(w, PageData{
		Artist:    *found,
		Relations: relations[id].DatesLocations,
	})
}

func SearchAPI(w http.ResponseWriter, r *http.Request) {
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

func FilterAPI(w http.ResponseWriter, r *http.Request) {
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

func NotFoundPage(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	message := "The artist you're looking for doesn't exist."
	if query != "" {
		message = "No artist found matching \"" + query + "\"."
	}
	RenderError(w, 404, "Artist not found", message)
}