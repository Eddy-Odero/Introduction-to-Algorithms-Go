package main

import (
	"strings"
	"testing"
)

// ─── ParseBanner Tests ────────────────────────────────────────────────────────

func TestParseBanner_ValidContent(t *testing.T) {
	// Build a minimal banner with just space and '!'
	var sb strings.Builder
	sb.WriteString("\n") // initial blank line

	// space (ASCII 32): 8 blank lines
	for i := 0; i < 8; i++ {
		sb.WriteString("      \n")
	}
	sb.WriteString("\n")

	// ! (ASCII 33): visible character
	exclamLines := []string{" _  ", "| | ", "| | ", "| | ", "|_| ", "(_) ", "    ", "    "}
	for _, l := range exclamLines {
		sb.WriteString(l + "\n")
	}
	sb.WriteString("\n")

	content := sb.String()
	charMap, err := ParseBanner(content)
	if err != nil {
		t.Fatalf("ParseBanner returned unexpected error: %v", err)
	}

	if _, ok := charMap[' ']; !ok {
		t.Error("charMap missing space character")
	}
	if _, ok := charMap['!']; !ok {
		t.Error("charMap missing '!' character")
	}
	if charMap['!'][0] != " _  " {
		t.Errorf("'!' line 0: got %q, want %q", charMap['!'][0], " _  ")
	}
}

func TestParseBanner_EmptyContent(t *testing.T) {
	_, err := ParseBanner("")
	if err == nil {
		t.Error("expected error for empty banner content, got nil")
	}
}

func TestParseBanner_CharacterCount(t *testing.T) {
	charMap := loadStandardForTest(t)
	// Standard printable ASCII: 32 to 126 inclusive = 95 characters
	if len(charMap) != 95 {
		t.Errorf("expected 95 characters, got %d", len(charMap))
	}
}

func TestParseBanner_EachCharHas8Lines(t *testing.T) {
	charMap := loadStandardForTest(t)
	for code := rune(32); code <= 126; code++ {
		lines, ok := charMap[code]
		if !ok {
			t.Errorf("missing character ASCII %d (%q)", code, string(code))
			continue
		}
		for i, line := range lines {
			// All lines of a character must have equal length
			if i > 0 && len(line) != len(lines[0]) {
				t.Errorf("char %q line %d length %d != line 0 length %d",
					string(code), i, len(line), len(lines[0]))
			}
		}
	}
}

// ─── LoadBanner Tests ─────────────────────────────────────────────────────────

func TestLoadBanner_Standard(t *testing.T) {
	charMap, err := LoadBanner("standard")
	if err != nil {
		t.Fatalf("LoadBanner(standard) error: %v", err)
	}
	if len(charMap) == 0 {
		t.Error("LoadBanner returned empty map for 'standard'")
	}
}

func TestLoadBanner_InvalidName(t *testing.T) {
	_, err := LoadBanner("invalid_banner")
	if err == nil {
		t.Error("expected error for unknown banner name, got nil")
	}
}

func TestLoadBanner_EmptyName(t *testing.T) {
	_, err := LoadBanner("")
	if err == nil {
		t.Error("expected error for empty banner name, got nil")
	}
}

// ─── RenderLine Tests ─────────────────────────────────────────────────────────

func TestRenderLine_SingleChar(t *testing.T) {
	charMap := loadStandardForTest(t)

	output, err := RenderLine("A", charMap)
	if err != nil {
		t.Fatalf("RenderLine error: %v", err)
	}

	rows := strings.Split(strings.TrimRight(output, "\n"), "\n")
	if len(rows) != 8 {
		t.Errorf("expected 8 rows for single char, got %d", len(rows))
	}
}

func TestRenderLine_EmptyString(t *testing.T) {
	charMap := loadStandardForTest(t)

	output, err := RenderLine("", charMap)
	if err != nil {
		t.Fatalf("RenderLine error: %v", err)
	}

	// Empty line should produce 8 newlines
	rows := strings.Split(output, "\n")
	// 8 lines each empty, split gives 9 entries (trailing \n)
	if len(rows) != 9 {
		t.Errorf("expected 9 parts after split, got %d", len(rows))
	}
}

func TestRenderLine_UnknownCharacter(t *testing.T) {
	charMap := loadStandardForTest(t)
	// Delete a character to simulate unknown char
	delete(charMap, 'A')

	_, err := RenderLine("A", charMap)
	if err == nil {
		t.Error("expected error for unknown character, got nil")
	}
}

func TestRenderLine_OutputHas8Rows(t *testing.T) {
	charMap := loadStandardForTest(t)

	for _, input := range []string{"a", "Z", "hello", "123"} {
		output, err := RenderLine(input, charMap)
		if err != nil {
			t.Fatalf("RenderLine(%q) error: %v", input, err)
		}
		rows := strings.Split(strings.TrimRight(output, "\n"), "\n")
		if len(rows) != 8 {
			t.Errorf("RenderLine(%q): expected 8 rows, got %d", input, len(rows))
		}
	}
}

// ─── PrintASCII Tests ─────────────────────────────────────────────────────────

func TestPrintASCII_EmptyString(t *testing.T) {
	result, err := PrintASCII("", "standard")
	if err != nil {
		t.Fatalf("PrintASCII(\"\") error: %v", err)
	}
	if result != "" {
		t.Errorf("PrintASCII(\"\") = %q, want empty string", result)
	}
}

func TestPrintASCII_SingleNewline(t *testing.T) {
	// go run . "\n" should output just a newline
	result, err := PrintASCII(`\n`, "standard")
	if err != nil {
		t.Fatalf("PrintASCII error: %v", err)
	}
	if result != "\n" {
		t.Errorf("PrintASCII(`\\n`) = %q, want single newline", result)
	}
}

func TestPrintASCII_HelloWorld(t *testing.T) {
	result, err := PrintASCII("Hello", "standard")
	if err != nil {
		t.Fatalf("PrintASCII(Hello) error: %v", err)
	}

	rows := strings.Split(strings.TrimRight(result, "\n"), "\n")
	if len(rows) != 8 {
		t.Errorf("PrintASCII(Hello): expected 8 rows, got %d", len(rows))
	}

	// All rows should have the same length (padded by character widths)
	firstLen := len(rows[0])
	for i, row := range rows {
		if len(row) != firstLen {
			t.Errorf("row %d length %d != row 0 length %d", i, len(row), firstLen)
		}
	}
}

func TestPrintASCII_NewlineInMiddle(t *testing.T) {
	// "Hello\nThere" should produce 16 rows (8 for Hello + 8 for There)
	result, err := PrintASCII(`Hello\nThere`, "standard")
	if err != nil {
		t.Fatalf("PrintASCII error: %v", err)
	}

	rows := strings.Split(strings.TrimRight(result, "\n"), "\n")
	if len(rows) != 16 {
		t.Errorf("Hello\\nThere: expected 16 rows, got %d", len(rows))
	}
}

func TestPrintASCII_DoubleNewline(t *testing.T) {
	// "Hello\n\nThere" should produce 8 + 1 (blank) + 8 = 17 rows
	result, err := PrintASCII(`Hello\n\nThere`, "standard")
	if err != nil {
		t.Fatalf("PrintASCII error: %v", err)
	}

	rows := strings.Split(strings.TrimRight(result, "\n"), "\n")
	if len(rows) != 17 {
		t.Errorf("Hello\\n\\nThere: expected 17 rows, got %d", len(rows))
	}
}

func TestPrintASCII_OnlyNewlines(t *testing.T) {
	// "\n" alone = one blank line
	result, err := PrintASCII(`\n`, "standard")
	if err != nil {
		t.Fatalf("PrintASCII error: %v", err)
	}
	if result != "\n" {
		t.Errorf(`PrintASCII(\n) = %q, want "\n"`, result)
	}
}

func TestPrintASCII_NumbersAndLetters(t *testing.T) {
	result, err := PrintASCII("1Hello", "standard")
	if err != nil {
		t.Fatalf("PrintASCII(1Hello) error: %v", err)
	}
	rows := strings.Split(strings.TrimRight(result, "\n"), "\n")
	if len(rows) != 8 {
		t.Errorf("1Hello: expected 8 rows, got %d", len(rows))
	}
}

func TestPrintASCII_SpecialCharacters(t *testing.T) {
	specials := []string{"!", "@", "#", "$", "%", "^", "&", "*", "(", ")"}
	for _, s := range specials {
		result, err := PrintASCII(s, "standard")
		if err != nil {
			t.Errorf("PrintASCII(%q) error: %v", s, err)
			continue
		}
		rows := strings.Split(strings.TrimRight(result, "\n"), "\n")
		if len(rows) != 8 {
			t.Errorf("PrintASCII(%q): expected 8 rows, got %d", s, len(rows))
		}
	}
}

func TestPrintASCII_Space(t *testing.T) {
	result, err := PrintASCII(" ", "standard")
	if err != nil {
		t.Fatalf("PrintASCII(space) error: %v", err)
	}
	rows := strings.Split(strings.TrimRight(result, "\n"), "\n")
	if len(rows) != 8 {
		t.Errorf("space: expected 8 rows, got %d", len(rows))
	}
}

func TestPrintASCII_InvalidBanner(t *testing.T) {
	_, err := PrintASCII("Hello", "nonexistent")
	if err == nil {
		t.Error("expected error for invalid banner, got nil")
	}
}

func TestPrintASCII_AllPrintableASCII(t *testing.T) {
	// Build a string with all printable ASCII characters
	var all strings.Builder
	for c := rune(32); c <= 126; c++ {
		all.WriteRune(c)
	}

	_, err := PrintASCII(all.String(), "standard")
	if err != nil {
		t.Errorf("PrintASCII(all printable ASCII) error: %v", err)
	}
}

// ─── Helper ───────────────────────────────────────────────────────────────────

func loadStandardForTest(t *testing.T) map[rune][8]string {
	t.Helper()
	charMap, err := LoadBanner("standard")
	if err != nil {
		t.Fatalf("Failed to load standard banner: %v", err)
	}
	return charMap
}
