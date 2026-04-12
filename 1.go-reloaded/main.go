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

	for i := 0; i < len(words); i++ {

		if words[i] == "(hex)" || words[i] == "(bin)" {
			if len(result) > 0 {
				prevValue := result[len(result)-1]

				var base int
				if words[i] == "(hex)" {
					base = 16
				} else if words[i] == "(bin)" {
					base = 2
				}

				decimal, err := strconv.ParseInt(prevValue, base, 64)
				if err == nil {
					result[len(result)-1] = fmt.Sprintf("%d", decimal)
				}
			}
			continue
		}

		result = append(result, words[i])
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