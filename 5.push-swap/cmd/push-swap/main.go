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

	normalized := utils.Normalize(numbers)

	a := stack.Stack{
		Data: normalized,
	}

	b := stack.Stack{}

	var operations []string

	switch len(numbers) {
	case 2:
		operations = sort.SortTwo(&a)

	case 3:
		operations = sort.SortThree(&a)

	case 4:
		operations = sort.SortFour(&a, &b)

	case 5:
		operations = sort.SortFive(&a, &b)

	default:
		operations = sort.RadixSort(&a, &b)
	}

	for _, op := range operations {
		fmt.Println(op)
	}
}