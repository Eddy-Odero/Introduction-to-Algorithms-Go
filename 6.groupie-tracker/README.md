# Groupie Tracker

A web application written in Go that fetches data from the [Groupie Trackers API](https://groupietrackers.herokuapp.com/api) and displays information about bands and artists, their members, concert dates, and tour locations.

## Features

- Browse all artists on a styled home page with cards and hover effects
- View individual artist pages with members and tour dates
- Live search — type in the search bar and get instant suggestions without a page reload
- Filter artists by founding year range and number of members
- Styled error pages for 404 and 500 responses

## Project Structure

```
groupie-tracker/
├── main.go                      # starts the server and registers routes
├── internal/
│   ├── api/
│   │   ├── api.go               # fetches and searches artist data from the API
│   │   └── api_test.go          # unit tests for search and filter logic
│   └── handlers/
│       ├── handlers.go          # HTTP handlers and template rendering
│       └── handlers_test.go     # unit tests for HTTP handlers
├── templates/
│   ├── index.html               # home page
│   ├── artist.html              # artist detail page
│   └── error.html               # error page (404, 500)
└── static/
    └── css/
        └── style.css            # styles
```

## API Used

The app fetches from four endpoints:

| Endpoint | What it contains |
|---|---|
| `/api/artists` | Band names, images, members, creation date, first album |
| `/api/locations` | Concert locations |
| `/api/dates` | Concert dates |
| `/api/relation` | Links artists to their dates and locations |

## Client-Server Events

The project implements two client-server events where the browser sends a request to the server and receives data back without a page reload:

**Live Search** — as the user types in the search bar, the browser sends a request to `/api/search?q=` and the server responds with matching artists in JSON.

**Filter** — when the user clicks the filter button, the browser sends a request to `/api/filter?minYear=&maxYear=&members=` and the server responds with filtered artists in JSON.

## Getting Started

**Requirements**
- Go 1.21 or higher
- Internet connection (fetches from external API)

**Run the server**
```bash
git clone https://github.com/Eddy-Odero/Introduction-to-Algorithms-Go.git
cd groupie-tracker
go run ./cmd/server/
```

Then open your browser at:
```
http://localhost:8080
```

**Run the tests**
```bash
go test ./...
```

## Allowed Packages

Only standard Go packages are used — no external dependencies.

## Authors

Built step by step as a learning project covering:
- JSON parsing and HTTP clients
- Go web servers with `net/http`
- HTML templates
- Client-server communication
- Unit testing with `net/http/httptest`