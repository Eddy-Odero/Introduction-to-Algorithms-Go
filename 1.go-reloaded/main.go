package main

import (
	"fmt"
	"os"
	"strings"
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
				fmt.Println("Found (hex), previous word is:", words[i-1])
			}
		}
	}
}