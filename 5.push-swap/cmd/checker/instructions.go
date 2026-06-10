package main

import (
	"fmt"

	"push-swap/stack"
)

func ExecuteInstruction(op string, a, b *stack.Stack) error {
	switch op {

	case "sa":
		stack.Sa(a)

	case "sb":
		stack.Sb(b)

	case "ss":
		stack.Ss(a, b)

	case "pa":
		stack.Pa(a, b)

	case "pb":
		stack.Pb(a, b)

	case "ra":
		stack.Ra(a)

	case "rb":
		stack.Rb(b)

	case "rr":
		stack.Rr(a, b)

	case "rra":
		stack.Rra(a)

	case "rrb":
		stack.Rrb(b)

	case "rrr":
		stack.Rrr(a, b)

	default:
		return fmt.Errorf("invalid instruction")
	}

	return nil
}