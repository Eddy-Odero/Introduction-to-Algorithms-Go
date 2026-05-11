package main

import "fmt"

func main(){
input := "((("
floor := 0
for _, i := range input{
	if i == '('{
		floor++
	}else if i == ')'{
		floor--
	}
	
}
fmt.Println(floor)
}