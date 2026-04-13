package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)
// Normalize tokens (fix "(up, 2)" split issue)

func normalize(words []string) []string {
	var res []string

	for i := 0; i < len(words); i++ {
		// merge "(up," + "2)" → "(up, 2)"
		if strings.HasSuffix(words[i], ",") && i+1 < len(words) && strings.HasSuffix(words[i+1], ")") {
			res = append(res, words[i]+" "+words[i+1])
			i++
		} else {
			res = append(res, words[i])
		}
	}
	return res
}
// Apply modifiers
func applyModifiers(words []string) []string {
	var result []string

	for i := 0; i < len(words); i++ {
		word := words[i]

		// ---- HEX / BIN ----
		if word == "(hex)" || word == "(bin)" {
			if len(result) > 0 {
				base := 16
				if word == "(bin)" {
					base = 2
				}

				prev := result[len(result)-1]
				if val, err := strconv.ParseInt(prev, base, 64); err == nil {
					result[len(result)-1] = fmt.Sprintf("%d", val)
				}
			}
			continue
		}

		// ---- TEXT MODIFIERS ----
		if strings.HasPrefix(word, "(up") ||
			strings.HasPrefix(word, "(low") ||
			strings.HasPrefix(word, "(cap") {

			count := 1

			// extract number if exists
			if strings.Contains(word, ",") {
				clean := strings.Trim(word, "()")
				parts := strings.Split(clean, ",")
				if len(parts) == 2 {
					if n, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
						count = n
					}
				}
			}

			for j := 0; j < count && len(result)-1-j >= 0; j++ {
				idx := len(result) - 1 - j
				w := result[idx]

				switch {
				case strings.HasPrefix(word, "(up"):
					result[idx] = strings.ToUpper(w)
				case strings.HasPrefix(word, "(low"):
					result[idx] = strings.ToLower(w)
				case strings.HasPrefix(word, "(cap"):
					if len(w) > 0 {
						result[idx] =
							strings.ToUpper(string(w[0])) +
							strings.ToLower(w[1:])
					}
				}
			}
			continue
		}

		result = append(result, word)
	}

	return result
}
// Fix "a" → "an"
func fixArticles(words []string) []string {
	for i := 0; i < len(words)-1; i++ {
		w := words[i]
		next := words[i+1]

		if strings.ToLower(w) == "a" && len(next) > 0 {
			first := unicode.ToLower(rune(next[0]))
			if strings.ContainsRune("aeiouh", first) {
				if w == "A" {
					words[i] = "An"
				} else {
					words[i] = "an"
				}
			}
		}
	}
	return words
}
// Fix punctuation + quotes
func fixPunctuation(words []string) string {
	punct := ".,!?;:"
	var result []string

	for i := 0; i < len(words); i++ {
		w := words[i]

		// punctuation attaches to previous word
		if len(w) > 0 && strings.ContainsRune(punct, rune(w[0])) {
			if len(result) > 0 {
				result[len(result)-1] += w
			} else {
				result = append(result, w)
			}
			continue
		}

		result = append(result, w)
	}

	// ---- FIX QUOTES ----
	var final []string
	for i := 0; i < len(result); i++ {
		if result[i] == "'" {
			j := i + 1
			for j < len(result) && result[j] != "'" {
				j++
			}
			if j < len(result) {
				inside := strings.Join(result[i+1:j], " ")
				final = append(final, "'"+inside+"'")
				i = j
				continue
			}
		}
		final = append(final, result[i])
	}

	return strings.Join(final, " ")
} 
// Full pipeline

func processLine(line string) string {
	words := strings.Fields(line)
	words = normalize(words)
	words = applyModifiers(words)
	words = fixArticles(words)
	return fixPunctuation(words)
}
// MAIN (file handling)
func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run . input.txt output.txt")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	lines := strings.Split(string(data), "\n")
	var output []string

	for _, line := range lines {
		output = append(output, processLine(line))
	}

	err = os.WriteFile(outputFile, []byte(strings.Join(output, "\n")), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
	}
}