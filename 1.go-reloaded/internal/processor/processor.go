package processor

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func normalize(words []string) []string {
	var res []string

	for i := 0; i < len(words); i++ {
		if strings.HasSuffix(words[i], ",") &&
			i+1 < len(words) &&
			strings.HasSuffix(words[i+1], ")") {

			res = append(res, words[i]+" "+words[i+1])
			i++
		} else {
			res = append(res, words[i])
		}
	}
	return res
}

func applyModifiers(words []string) []string {
	var result []string

	for i := 0; i < len(words); i++ {
		word := words[i]

		// HEX / BIN
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

		// TEXT MODIFIERS
		if strings.HasPrefix(word, "(up") ||
			strings.HasPrefix(word, "(low") ||
			strings.HasPrefix(word, "(cap") {

			count := 1

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

func fixPunctuation(words []string) string {
	punct := ".,!?;:"

	var out []string

	for i := 0; i < len(words); i++ {
		w := words[i]

		for len(w) > 0 && strings.ContainsRune(punct, rune(w[0])) {
			if len(out) > 0 {
				out[len(out)-1] += string(w[0])
			} else {
				out = append(out, string(w[0]))
			}
			w = w[1:]
		}

		if len(w) > 0 {
			out = append(out, w)
		}
	}
	if len(out) == 0 {
		return ""
	}

	result := out[0]

	for i := 1; i < len(out); i++ {
    curr := out[i]
    prev := out[i-1] // ← use out[i-1], not last char of result

    // punctuation attaches to previous word (no space)
    if strings.ContainsRune(punct, rune(curr[0])) {
        result += curr
        continue
    }

    // previous token was pure punctuation (e.g. "...") → no space
    isPurePunct := true
    for _, ch := range prev {
        if !strings.ContainsRune(punct, ch) {
            isPurePunct = false
            break
        }
    }
    if isPurePunct {
        result += curr
        continue
    }

    result += " " + curr
}
	for {
		start := strings.Index(result, "' ")
		if start == -1 {
			break
		}

		end := strings.Index(result[start+2:], " '")
		if end == -1 {
			break
		}

		end += start + 2

		inside := strings.TrimSpace(result[start+2 : end])
		result = result[:start] + "'" + inside + "'" + result[end+2:]
	}

	return result
}

func ProcessLine(line string) string {
	words := strings.Fields(line)
	words = normalize(words)
	words = applyModifiers(words)
	words = fixArticles(words)
	return fixPunctuation(words)
}