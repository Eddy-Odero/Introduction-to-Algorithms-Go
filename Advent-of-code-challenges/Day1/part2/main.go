package main

import "fmt"

func main() {
    input := "()())"

    floor := 0

    for pos, ch := range input {
        if ch == '(' {
            floor++
        } else if ch == ')' {
            floor--
        }

        if floor == -1 {
            fmt.Println("First basement position:", pos+1)
            break
        }
    }
}