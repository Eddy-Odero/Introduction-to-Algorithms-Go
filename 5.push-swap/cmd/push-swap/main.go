package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	args := os.Args[1:]

	var values []int

	for _, arg := range args {
		fields := strings.Fields(arg)

		for _, field := range fields {
			n, err := strconv.Atoi(field)

			if err != nil {
				fmt.Println("Error")
				return
			}

			values = append(values, n)
		}
	}

	fmt.Println(values)
}