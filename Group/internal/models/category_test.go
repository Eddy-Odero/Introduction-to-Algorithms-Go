package models_test

import (
	"testing"

	"forum/internal/models"
)

func TestCategoriesSeeded(t *testing.T) {
	conn := openTestDB(t)
	cats := &models.CategoryStore{DB: conn}

	all, err := cats.All()
	if err != nil {
		t.Fatalf("All: %v", err)
	}

	want := map[string]bool{"general": true, "technology": true, "gaming": true, "random": true}
	if len(all) != len(want) {
		t.Fatalf("got %d categories, want %d", len(all), len(want))
	}
	for _, c := range all {
		if !want[c.Name] {
			t.Errorf("unexpected category %q", c.Name)
		}
	}
}

func TestIDsByNamesSkipsUnknown(t *testing.T) {
	conn := openTestDB(t)
	cats := &models.CategoryStore{DB: conn}

	ids, err := cats.IDsByNames([]string{"technology", "not-a-real-category", "gaming"})
	if err != nil {
		t.Fatalf("IDsByNames: %v", err)
	}
	if len(ids) != 2 {
		t.Fatalf("got %d ids, want 2 (unknown category should be skipped)", len(ids))
	}
}
