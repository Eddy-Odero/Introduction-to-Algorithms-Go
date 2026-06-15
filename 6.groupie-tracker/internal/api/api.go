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

type Relations struct {
    ID             int                 `json:"id"`
    DatesLocations map[string][]string `json:"datesLocations"`
}

type RelationsResponse struct {
    Index []Relations `json:"index"`
}

func GetRelations() (map[int]Relations, error) {
    resp, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
    if err != nil {
        return nil, fmt.Errorf("fetch failed: %w", err)
    }
    defer resp.Body.Close()

    var result RelationsResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("decode failed: %w", err)
    }

    // turn the slice into a map keyed by artist ID for easy lookup
    relMap := make(map[int]Relations)
    for _, r := range result.Index {
        relMap[r.ID] = r
    }

    return relMap, nil
}
func FilterArtists(artists []Artist, minYear, maxYear, members int) []Artist {
	var results []Artist

	for _, a := range artists {
		// check year range
		if minYear > 0 && a.CreationDate < minYear {
			continue
		}
		if maxYear > 0 && a.CreationDate > maxYear {
			continue
		}
		// check members count
		if members > 0 && len(a.Members) != members {
			continue
		}
		results = append(results, a)
	}

	return results
}