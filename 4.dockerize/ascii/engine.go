package ascii
import "strings"

func BuildFont(lines []string) map[rune][]string {
	font := make(map[rune][]string)

	char := ' '
	row := 0

	for i := 0; i < len(lines); i++ {

		if lines[i] == "" {
			continue
		}

		if row%8 == 0 && row != 0 {
			char++
		}

		font[char] = append(font[char], lines[i])
		row++
	}

	return font
}
func Generate(text string, font map[rune][]string) string {
	if text == "" {
		return ""
	}

	result := make([]string, 8)

	for _, char := range text {
		asciiChar, ok := font[char]
		if !ok {
			continue
		}

		for i := 0; i < 8; i++ {
			result[i] += asciiChar[i]
		}
	}

	return strings.Join(result, "\n")
}