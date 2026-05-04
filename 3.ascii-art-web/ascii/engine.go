package ascii

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
	result := make([]string, 8)

	for _, c := range text {

		ascii, ok := font[c]
		if !ok {
			continue
		}

		for i := 0; i < 8; i++ {
			result[i] += ascii[i]
		}
	}

	out := ""
	for _, line := range result {
		out += line + "\n"
	}

	return out
}