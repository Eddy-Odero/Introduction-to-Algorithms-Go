package handlers

import (
    "encoding/json"
    "groupie-tracker/internal/api"
    "html/template"
    "net/http"
    "path/filepath"
    "runtime"
    "strconv"
    "strings"
)
var client = &api.Client{}

type PageData struct {
	Artist    api.Artist
	Relations map[string][]string
}

type ErrorData struct {
	Code    int
	Title   string
	Message string
}

func projectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "..")
}

func RenderError(w http.ResponseWriter, code int, title, message string) {
	tmpl, err := template.ParseFiles(filepath.Join(projectRoot(), "templates", "error.html"))
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
	artists, err := client.GetArtists()
	if err != nil {
		RenderError(w, 500, "Something went wrong", "We couldn't load the artists. Please try again.")
		return
	}

	tmpl, err := template.ParseFiles(filepath.Join(projectRoot(), "templates", "index.html"))
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

	// channels to collect results from both goroutines
	artistsCh := make(chan []api.Artist, 1)
	relationsCh := make(chan map[int]api.Relations, 1)
	errCh := make(chan error, 2)

	// fetch artists in a goroutine
	go func() {
		artists, err := client.GetArtists()
		if err != nil {
			errCh <- err
			return
		}
		artistsCh <- artists
	}()

	// fetch relations in a goroutine — runs at the same time as above
	go func() {
		relations, err := client.GetRelations()
		if err != nil {
			errCh <- err
			return
		}
		relationsCh <- relations
	}()

	// wait for both to finish
	var artists []api.Artist
	var relations map[int]api.Relations

	for i := 0; i < 2; i++ {
		select {
		case a := <-artistsCh:
			artists = a
		case rel := <-relationsCh:
			relations = rel
		case err := <-errCh:
			RenderError(w, 500, "Something went wrong", err.Error())
			return
		}
	}

	// find the artist
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

	tmpl, err := template.ParseFiles(filepath.Join(projectRoot(), "templates", "artist.html"))
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

	artists, err := client.GetArtists()
	if err != nil {
		http.Error(w, "Failed to load artists", http.StatusInternalServerError)
		return
	}

	results := client.SearchArtists(query, artists)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func FilterAPI(w http.ResponseWriter, r *http.Request) {
	minYear, _ := strconv.Atoi(r.URL.Query().Get("minYear"))
	maxYear, _ := strconv.Atoi(r.URL.Query().Get("maxYear"))
	members, _ := strconv.Atoi(r.URL.Query().Get("members"))

	artists, err := client.GetArtists()
	if err != nil {
		http.Error(w, "Failed to load artists", http.StatusInternalServerError)
		return
	}

	results := client.FilterArtists(artists, minYear, maxYear, members)

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