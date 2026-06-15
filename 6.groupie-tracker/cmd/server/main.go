package main

import (
    "fmt"
    "html/template"
    "net/http"
    "strconv"
    "strings"
	"encoding/json"
    "groupie-tracker/internal/api"
)

func home(w http.ResponseWriter, r *http.Request) {
    artists, err := api.GetArtists()
    if err != nil {
        http.Error(w, "Failed to load artists", http.StatusInternalServerError)
        return
    }

    tmpl, err := template.ParseFiles("templates/index.html")
    if err != nil {
        http.Error(w, "Failed to load template", http.StatusInternalServerError)
        return
    }

    tmpl.Execute(w, artists)
}

func artistPage(w http.ResponseWriter, r *http.Request) {
    // extract the ID from the URL e.g. "/artist/3" → "3"
    idStr := strings.TrimPrefix(r.URL.Path, "/artist/")
    id, err := strconv.Atoi(idStr)
    if err != nil || id < 1 {
        http.Error(w, "Artist not found", http.StatusNotFound)
        return
    }

    artists, err := api.GetArtists()
    if err != nil {
        http.Error(w, "Failed to load artists", http.StatusInternalServerError)
        return
    }

    // find the artist with matching ID
    var found *api.Artist
    for _, a := range artists {
        if a.ID == id {
            found = &a
            break
        }
    }

    if found == nil {
        http.Error(w, "Artist not found", http.StatusNotFound)
        return
    }

    tmpl, err := template.ParseFiles("templates/artist.html")
    if err != nil {
        http.Error(w, "Failed to load template", http.StatusInternalServerError)
        return
    }

    tmpl.Execute(w, found)
}


func searchAPI(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("q")

    artists, err := api.GetArtists()
    if err != nil {
        http.Error(w, "Failed to load artists", http.StatusInternalServerError)
        return
    }

    results := api.SearchArtists(query, artists)

    // respond with JSON
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(results)
}

func main() {
    // serve static files
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    http.HandleFunc("/api/search", searchAPI)
    http.HandleFunc("/artist/", artistPage)
    http.HandleFunc("/", home)

    fmt.Println("Server running at http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}