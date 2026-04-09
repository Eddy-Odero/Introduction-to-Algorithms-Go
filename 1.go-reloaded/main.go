package main

import (
	"fmt"
	"os"
)

func main() {
	data, err := os.ReadFile("quiz.txt")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf(string(data))
}
