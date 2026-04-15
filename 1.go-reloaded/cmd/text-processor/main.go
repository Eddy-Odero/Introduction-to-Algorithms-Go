package main

import (
	"fmt"
	"os"
	"strings"

	"text-processor/internal/processor"
)

func main() {
	// check arguments
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run ./cmd/text-processor input.txt output.txt")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	// read input file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}

	lines := strings.Split(string(data), "\n")

	var result []string

	for _, line := range lines {
		result = append(result, processor.ProcessLine(line))
	}

	// write output file
	err = os.WriteFile(outputFile, []byte(strings.Join(result, "\n")), 0644)
	if err != nil {
		fmt.Println("Error writing output file:", err)
	}
}
// run (go run ./cmd/text-processor testdata/sample.txt testdata/result.txt) to test the project