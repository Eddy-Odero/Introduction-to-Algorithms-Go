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

	for i, word := range words {
		if word == "(hex)" {
			if i > 0 {
				hexValue := words[i-1]

				decimal, err := strconv.ParseInt(hexValue, 16, 64)
				if err != nil {
					fmt.Println("Conversion error:", err)
					continue
				}

				fmt.Println("Hex:", hexValue, "→ Decimal:", decimal)
			}
		}
	}
}