package main

import (
	"fmt"
	"os"

	"Tetris-optimizer/internal/parser"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("ERROR")
		return
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("ERROR")
		return
	}

	pieces, err := parser.Parse(string(data))
	if err != nil {
		fmt.Println("ERROR")
		return
	}
for i, piece := range pieces {
	fmt.Println("Piece", i)

	for _, row := range piece.Raw {
		fmt.Println(row)
	}

	fmt.Println()
}
	fmt.Println("Pieces found:", len(pieces))
	fmt.Printf("%q\n", string(data))
}