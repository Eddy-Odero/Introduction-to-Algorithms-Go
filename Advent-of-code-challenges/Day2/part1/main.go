package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	// Read the whole file
	data, err := os.ReadFile("input.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Convert the file to a string and split into lines
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")

	totalPaper := 0

	for _, line := range lines {
		var l, w, h int

		// Read dimensions like "2x3x4"
		fmt.Sscanf(line, "%dx%dx%d", &l, &w, &h)

		// Calculate the area of each side
		side1 := l * w
		side2 := w * h
		side3 := h * l

		// Surface area
		surface := 2*side1 + 2*side2 + 2*side3

		// Find the smallest side
		smallest := side1
		if side2 < smallest {
			smallest = side2
		}
		if side3 < smallest {
			smallest = side3
		}

		totalPaper += surface + smallest
	}

	fmt.Println("Total wrapping paper:", totalPaper)
}