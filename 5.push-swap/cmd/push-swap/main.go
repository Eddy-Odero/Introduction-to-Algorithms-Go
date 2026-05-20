package main

import (
	"fmt"
	"os"

	"push-swap/parse"
)

func main() {
	numbers, err := parse.ParseArgs(os.Args[1:])

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error")
		return
	}

	fmt.Println(numbers)
}