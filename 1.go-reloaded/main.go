package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func processLine(line string) string {
	words := strings.Fields(line)
	var result []string

	bases := map[string]int{
		"(hex)": 16,
		"(bin)": 2,
	}

	for i := 0; i < len(words); i++ {

		word := words[i]

		// ---- Number conversions ----
		if base, ok := bases[word]; ok {
			if len(result) > 0 {
				prev := result[len(result)-1]

				decimal, err := strconv.ParseInt(prev, base, 64)
				if err == nil {
					result[len(result)-1] = fmt.Sprintf("%d", decimal)
				}
			}
			continue
		}

		// ---- Text transformations ----
		if len(result) > 0 {
			prev := result[len(result)-1]

			switch word {
			case "(up)":
				result[len(result)-1] = strings.ToUpper(prev)
				continue

			case "(low)":
				result[len(result)-1] = strings.ToLower(prev)
				continue

			case "(cap)":
				if len(prev) > 0 {
					result[len(result)-1] =
						strings.ToUpper(string(prev[0])) +
						strings.ToLower(prev[1:])
				}
				continue
			}
		}

		// Normal word
		result = append(result, word)
	}

	return strings.Join(result, " ")
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