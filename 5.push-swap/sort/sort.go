package sort

import "push-swap/stack"

func SortTwo(a *stack.Stack) []string {
	if len(a.Data) != 2 {
		return nil
	}

	if a.Data[0] > a.Data[1] {
		stack.Sa(a)
		return []string{"sa"}
	}

	return nil
}