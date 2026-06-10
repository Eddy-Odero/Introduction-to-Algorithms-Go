package main

import (
	"fmt"
	"os"

	"push-swap/parse"
	"push-swap/sort"
	"push-swap/stack"
	"push-swap/utils"
)

func main() {
	numbers, err := parse.ParseArgs(os.Args[1:])

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error")
		return
	}

	if utils.IsSorted(numbers) {
		return
	}

	a := stack.Stack{
		Data: numbers,
	}

	var operations []string

	if len(numbers) == 2 {
		operations = sort.SortTwo(&a)
	}

	for _, op := range operations {
		fmt.Println(op)
	}
}