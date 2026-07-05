package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("ERROR: invalid data format, wrong number of arguments")
		return
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("error reading file:")
		return
	}

	content := string(data)
	lines := strings.Split(content , "\n")

	for i, line := range lines {
		fmt.Println( "line",i, line)
	}
}