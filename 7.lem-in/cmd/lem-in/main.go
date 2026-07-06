package main

import (
	"fmt"
	"os"
	"strings"

	"lem-in/internal/graph"
	"lem-in/internal/parser"
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

	fmt.Println("Ants:", colony.NumAnts)
	fmt.Println("Start:", colony.Start.Name, colony.Start.X, colony.Start.Y)
	fmt.Println("End:", colony.End.Name, colony.End.X, colony.End.Y)
	fmt.Println("Total rooms:", len(colony.Rooms))
	for name, room := range colony.Rooms {
		linkNames := []string{}
		for _, l := range room.Links {
			linkNames = append(linkNames, l.Name)
		}
		fmt.Printf("  %s -> %v\n", name, linkNames)
	}

	path := graph.FindShortestPath(colony)
	if path == nil {
		fmt.Println("No path found")
	} else {
		fmt.Print("Shortest path: ")
		for _, r := range path {
			fmt.Print(r.Name, " ")
		}
		fmt.Println()
	}
}