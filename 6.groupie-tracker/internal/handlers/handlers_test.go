package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	
)



func TestHomeReturns200(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	Home(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET / expected 200, got %d", w.Code)
	}
}

func TestArtistPageValidID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/artist/1", nil)
	req.URL.Path = "/artist/1"
	w := httptest.NewRecorder()

	ArtistPage(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /artist/1 expected 200, got %d", w.Code)
	}
}

func TestArtistPageInvalidID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/artist/abc", nil)
	req.URL.Path = "/artist/abc"
	w := httptest.NewRecorder()

	ArtistPage(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("GET /artist/abc expected 404, got %d", w.Code)
	}
}

func TestArtistPageNotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/artist/999", nil)
	req.URL.Path = "/artist/999"
	w := httptest.NewRecorder()

	ArtistPage(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("GET /artist/999 expected 404, got %d", w.Code)
	}
}

func TestSearchAPIReturnsJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/search?q=queen", nil)
	w := httptest.NewRecorder()

	SearchAPI(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /api/search?q=queen expected 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}
}

func TestSearchAPIEmptyQuery(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/search?q=", nil)
	w := httptest.NewRecorder()

	SearchAPI(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /api/search?q= expected 200, got %d", w.Code)
	}
}

func TestFilterAPIReturnsJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/filter?minYear=1960&maxYear=1990", nil)
	w := httptest.NewRecorder()

	FilterAPI(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /api/filter expected 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}
}

func TestNotFoundPage(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/notfound?q=xyz", nil)
	w := httptest.NewRecorder()

	NotFoundPage(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("GET /notfound expected 404, got %d", w.Code)
	}
}
