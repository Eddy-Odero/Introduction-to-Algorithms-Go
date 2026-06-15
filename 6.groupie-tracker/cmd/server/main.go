package main

import (
    "fmt"
    "groupie-tracker/internal/api"
)

func main() {
    artists, err := api.GetArtists()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    for _, a := range artists {
        fmt.Printf("ID:%d  Name:%-25s  Founded:%d  Members:%v\n",
            a.ID, a.Name, a.CreationDate, a.Members)
    }
}