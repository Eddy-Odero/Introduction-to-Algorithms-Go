package stack

import "fmt"

type Stack struct {
	Data []int
}

func (s Stack) Print(name string) {
	fmt.Println(name + ":")

	for _, v := range s.Data {
		fmt.Println(v)
	}

	fmt.Println()
}