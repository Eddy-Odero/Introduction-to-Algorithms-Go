package main

import (
	"bufio"
	"fmt"
	"os"

	"push-swap/parse"
	"push-swap/stack"
	"push-swap/utils"
)

func main() {
	numbers, err := parse.ParseArgs(os.Args[1:])

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error")
		return
	}

	a := stack.Stack{Data: numbers}
	b := stack.Stack{}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		instruction := scanner.Text()

		err := ExecuteInstruction(instruction, &a, &b)

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error")
			return
		}
	}

	if utils.IsSorted(a.Data) && len(b.Data) == 0 {
		fmt.Println("OK")
	} else {
		fmt.Println("KO")
	}
}