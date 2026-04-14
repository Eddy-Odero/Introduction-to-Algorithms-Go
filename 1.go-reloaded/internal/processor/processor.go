package processor

import "strings"

func ProcessLine(line string) string {
	words := strings.Fields(line)

	var result []string
	punct := ".,!?;:"

	for _, w := range words {
		// if punctuation → attach to previous word
		if len(w) > 0 && strings.ContainsRune(punct, rune(w[0])) {
			if len(result) > 0 {
				result[len(result)-1] += w
			} else {
				result = append(result, w)
			}
		} else {
			result = append(result, w)
		}
	}

	return strings.Join(result, " ")
}