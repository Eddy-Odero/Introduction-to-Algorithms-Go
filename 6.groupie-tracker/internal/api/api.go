package api

import (
    "encoding/json"
    "fmt"
    "net/http"
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