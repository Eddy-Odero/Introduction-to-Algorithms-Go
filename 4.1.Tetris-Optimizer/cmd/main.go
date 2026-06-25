package main

import ("fmt"
"os"
)

func main(){
	if len(os.Args) != 2{
		 fmt.Println("Error provide file path")
	}
	file := os.Args[1]
data, err := os.ReadFile(file)
if err != nil{
	fmt.Println("error reading file", err)
}
fmt.Println(string(data))
}