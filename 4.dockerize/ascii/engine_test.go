package ascii

import "testing"

func TestLoadBanner(t *testing.T) {
	_, err := LoadBanner("../banners/standard.txt")
	if err != nil {
		t.Errorf("Failed to load banner: %v", err)
	}
}

func TestLoadBanner_InvalidFile(t *testing.T) {
	_, err := LoadBanner("invalid.txt")
	if err == nil {
		t.Errorf("Expected error for invalid file")
	}
}

func TestGenerate_SingleChar(t *testing.T) {
	lines, _ := LoadBanner("../banners/standard.txt")
	font := BuildFont(lines)

	result := Generate("A", font)

	if result == "" {
		t.Errorf("Expected output, got empty string")
	}
}

func TestGenerate_EmptyInput(t *testing.T) {
	lines, _ := LoadBanner("../banners/standard.txt")
	font := BuildFont(lines)

	result := Generate("", font)

	if result != "" {
		t.Errorf("Expected empty result for empty input")
	}
}

func TestGenerate_MultipleChars(t *testing.T) {
	lines, _ := LoadBanner("../banners/standard.txt")
	font := BuildFont(lines)

	result := Generate("AB", font)

	if len(result) == 0 {
		t.Errorf("Expected output for multiple characters")
	}
}

func TestGenerate_UnsupportedChar(t *testing.T) {
	lines, _ := LoadBanner("../banners/standard.txt")
	font := BuildFont(lines)

	result := Generate("A☺", font)

	if result == "" {
		t.Errorf("Expected partial output, got empty")
	}
}