package main

import (
	"fmt"
	"os"

	"push-swap/parse"
	"push-swap/stack"
)
func main() {
	numbers, err := parse.ParseArgs(os.Args[1:])

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error")
		return
	}

	a := stack.Stack{
		Data: numbers,
	}

	a.Print("A")
}