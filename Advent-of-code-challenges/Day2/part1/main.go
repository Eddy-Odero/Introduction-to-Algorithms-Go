package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	totalPaper := 0

	for scanner.Scan() {
		var l, w, h int

		// Read dimensions in the form LxWxH
		fmt.Sscanf(scanner.Text(), "%dx%dx%d", &l, &w, &h)

		// Areas of the three sides
		side1 := l * w
		side2 := w * h
		side3 := h * l

		// Surface area
		surfaceArea := 2*side1 + 2*side2 + 2*side3

		// Find the smallest side
		smallest := side1
		if side2 < smallest {
			smallest = side2
		}
		if side3 < smallest {
			smallest = side3
		}

		// Total paper for this present
		totalPaper += surfaceArea + smallest
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	fmt.Println("Total wrapping paper needed:", totalPaper)
}