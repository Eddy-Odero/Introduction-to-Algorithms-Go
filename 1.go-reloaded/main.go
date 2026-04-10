package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"
)

func main() {
	data, err := os.ReadFile("test.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	text := string(data)
words := strings.Split(text, " ")
	var result []string

	for i := 0; i < len(words); i++ {

		if words[i] == "(hex)" {
			// Convert previous word
			if len(result) > 0 {
				hexValue := result[len(result)-1]

				decimal, err := strconv.ParseInt(hexValue, 16, 64)
				if err == nil {
					// Replace last word with decimal
					result[len(result)-1] = fmt.Sprintf("%d", decimal)
				}
			}
			// Skip adding "(hex)"
			continue
		}

		// Normal word → add to result
		result = append(result, words[i])
	}

	// Join back into sentence
	final := strings.Join(result, " ")
	fmt.Println(final)
}