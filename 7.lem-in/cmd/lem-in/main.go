package main

import (
	"fmt"
	"os"
	"strings"

	"lem-in/internal/graph"
	"lem-in/internal/parser"
	"lem-in/internal/simulator"
	"lem-in/internal/solver"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("ERROR: invalid data format, wrong number of arguments")
		return
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("ERROR: invalid data format, cannot read file")
		return
	}

	lines := strings.Split(string(data), "\n")

	colony, err := parser.Parse(lines)
	if err != nil {
		fmt.Println(err)
		return
	}

	allPaths := graph.FindAllPaths(colony)
	if len(allPaths) == 0 {
		fmt.Println("ERROR: invalid data format, no path found")
		return
	}

	displayLines := colony.RawLines
	if len(displayLines) > 0 && displayLines[len(displayLines)-1] == "" {
		displayLines = displayLines[:len(displayLines)-1]
	}
	fmt.Println(strings.Join(displayLines, "\n"))
	fmt.Println()

	assignments := solver.AssignAnts(colony.NumAnts, allPaths)

	// 3. Simulate turn-by-turn movement and print each turn.
	turns := simulator.Simulate(assignments)
	for _, t := range turns {
		fmt.Println(t)
	}
}