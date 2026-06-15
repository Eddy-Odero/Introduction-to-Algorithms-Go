package api

import (
    "encoding/json"
    "fmt"
    "net/http"
	"strings"
)

// Artist matches the shape of one object in the API response
type Artist struct {
    ID           int      `json:"id"`
    Name         string   `json:"name"`
    Image        string   `json:"image"`
    Members      []string `json:"members"`
    CreationDate int      `json:"creationDate"`
    FirstAlbum   string   `json:"firstAlbum"`
}

func GetArtists() ([]Artist, error) {
    resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
    if err != nil {
        return nil, fmt.Errorf("fetch failed: %w", err)
    }
    defer resp.Body.Close()

    var artists []Artist
    if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
        return nil, fmt.Errorf("decode failed: %w", err)
    }

    return artists, nil
}

func SearchArtists(query string, artists []Artist) []Artist {
    if query == "" {
        return nil
    }

    query = strings.ToLower(query)
    var results []Artist

    for _, a := range artists {
        if strings.Contains(strings.ToLower(a.Name), query) {
            results = append(results, a)
            continue
        }
        for _, m := range a.Members {
            if strings.Contains(strings.ToLower(m), query) {
                results = append(results, a)
                break
            }
        }
    }

    return results
}