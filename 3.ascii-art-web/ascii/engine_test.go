package ascii

import "testing"

func TestLoadBanner(t *testing.T) {
	_, err := LoadBanner("../banners/standard.txt")
	if err != nil {
		t.Errorf("Failed to load banner: %v", err)
	}
}

func TestBuildFont(t *testing.T) {
	lines, err := LoadBanner("../banners/standard.txt")
	if err != nil {
		t.Fatalf("Could not load banner: %v", err)
	}

	font := BuildFont(lines)

	if len(font) == 0 {
		t.Errorf("Font map is empty")
	}
}

func TestGenerate(t *testing.T) {
	lines, _ := LoadBanner("../banners/standard.txt")
	font := BuildFont(lines)

	result := Generate("A", font)

	if result == "" {
		t.Errorf("Generate returned empty result")
	}
}