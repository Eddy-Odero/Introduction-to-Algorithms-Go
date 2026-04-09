package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	data, err := os.ReadFile("quiz.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	content := string(data)
lines := strings.Split(content, "\n")

for i, line := range lines {
	if strings.Contains(line, "Go") {
		fmt.Println("Found in line", i+1, ":", line)
	}
}
}