package main

import (
	"fmt"
	"os"
)

func main() {
	data, err := os.ReadFile("testdata/sample.txt")
	if err != nil {
		fmt.Println("error reading file", err)
		return
	}

	fmt.Println(string(data))
}
