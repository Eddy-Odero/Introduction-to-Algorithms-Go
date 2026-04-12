package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func processLine(line string) string {
	words := strings.Fields(line)

	// -------------------------
	// PASS 1: Apply modifiers
	// -------------------------
	var result []string

	for i := 0; i < len(words); i++ {
		word := words[i]

		// ---------- (hex), (bin) ----------
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

		// ---------- (up), (low), (cap) with optional number ----------
		if strings.HasPrefix(word, "(up") ||
			strings.HasPrefix(word, "(low") ||
			strings.HasPrefix(word, "(cap") {

			// default = 1 word
			count := 1

			// check for (up, 2)
			if strings.Contains(word, ",") {
				clean := strings.Trim(word, "()")
				parts := strings.Split(clean, ",")
				if len(parts) == 2 {
					if n, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
						count = n
					}
				}
			}

			// apply to last N words
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

	// -------------------------
	// PASS 2: a → an
	// -------------------------
	for i := 0; i < len(result)-1; i++ {
		if strings.ToLower(result[i]) == "a" {
			next := strings.ToLower(result[i+1])
			if len(next) > 0 && strings.ContainsAny(string(next[0]), "aeiou") {
				result[i] = "an"
			}
		}
	}

	// -------------------------
	// PASS 3: punctuation
	// -------------------------
	var final []string
	punct := ".,!?;:"

	for i := 0; i < len(result); i++ {
		w := result[i]

		// handle punctuation
		if len(w) > 0 && strings.ContainsAny(string(w[0]), punct) {
			if len(final) > 0 {
				final[len(final)-1] += w
			} else {
				final = append(final, w)
			}
			continue
		}

		final = append(final, w)
	}

	// handle quotes '
	for i := 0; i < len(final); i++ {
		if final[i] == "'" && i+2 < len(final) {
			j := i + 2
			for j < len(final) && final[j] != "'" {
				j++
			}
			if j < len(final) {
				// join inside quotes
				joined := strings.Join(final[i+1:j], " ")
				final = append(final[:i], append([]string{"'" + joined + "'"}, final[j+1:]...)...)
			}
		}
	}

	return strings.Join(final, " ")
}

func main() {
	data, err := os.ReadFile("test.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		processed := processLine(line)
		fmt.Println(processed)
	}
}