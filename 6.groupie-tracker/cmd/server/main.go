package main

import (
    "fmt"
    "net/http"
    "groupie-tracker/internal/api"
)

func home(w http.ResponseWriter, r *http.Request) {
    artists, err := api.GetArtists()
    if err != nil {
        http.Error(w, "Failed to load artists", http.StatusInternalServerError)
        return
    }

    for _, a := range artists {
        fmt.Fprintf(w, "ID:%d  Name:%-25s  Founded:%d\n", a.ID, a.Name, a.CreationDate)
    }
}

func main() {
    http.HandleFunc("/", home)

    fmt.Println("Server running at http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}