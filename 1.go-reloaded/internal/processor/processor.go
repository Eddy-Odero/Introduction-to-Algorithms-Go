package processor

import (
	"strconv"
	"strings"
)

func ProcessLine(line string) string {
	words := strings.Fields(line)

	var result []string

	for i := 0; i < len(words); i++ {
		w := words[i]

		// -------------------------
		// (up)
		// -------------------------
		if w == "(up)" {
			for j := 0; j < 1 && len(result)-1-j >= 0; j++ {
				idx := len(result) - 1 - j
				result[idx] = strings.ToUpper(result[idx])
			}
			continue
		}

		// -------------------------
		// (up, N) → TWO TOKENS: "(up," + "2)"
		// -------------------------
		if w == "(up," && i+1 < len(words) {

			numStr := strings.TrimRight(words[i+1], ")")

			count := 1
			if n, err := strconv.Atoi(numStr); err == nil {
				count = n
			}

			for j := 0; j < count && len(result)-1-j >= 0; j++ {
				idx := len(result) - 1 - j
				result[idx] = strings.ToUpper(result[idx])
			}

			i++ // skip number token
			continue
		}

		// -------------------------
		// punctuation handling
		// -------------------------
		if len(w) == 1 && strings.Contains(".,!?;:", w) {
			if len(result) > 0 {
				result[len(result)-1] += w
			}
			continue
		}

		// normal word
		result = append(result, w)
	}

	return strings.Join(result, " ")
}