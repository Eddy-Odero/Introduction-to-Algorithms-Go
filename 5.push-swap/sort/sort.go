package sort

import ("push-swap/stack"
"push-swap/utils"
)

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

func SortThree(a *stack.Stack) []string {
	if len(a.Data) != 3 {
		return nil
	}

	x := a.Data[0]
	y := a.Data[1]
	z := a.Data[2]

	// 2 1 3
	if x > y && y < z && x < z {
		stack.Sa(a)
		return []string{"sa"}
	}

	// 3 1 2
	if x > y && y < z && x > z {
		stack.Ra(a)
		return []string{"ra"}
	}

	// 1 3 2
	if x < y && y > z && x < z {
		stack.Sa(a)
		stack.Ra(a)
		return []string{"sa", "ra"}
	}

	// 3 2 1
	if x > y && y > z {
		stack.Sa(a)
		stack.Rra(a)
		return []string{"sa", "rra"}
	}

	// 2 3 1
	if x < y && y > z && x > z {
		stack.Rra(a)
		return []string{"rra"}
	}

	return nil
}

func moveMinToTop(a *stack.Stack, ops *[]string) {
	min := utils.Min(a.Data)
	index := utils.IndexOf(a.Data, min)

	if index <= len(a.Data)/2 {
		for index > 0 {
			stack.Ra(a)
			*ops = append(*ops, "ra")
			index--
		}
	} else {
		for index < len(a.Data) {
			stack.Rra(a)
			*ops = append(*ops, "rra")

			index++

			if a.Data[0] == min {
				break
			}
		}
	}
}
func SortFour(a, b *stack.Stack) []string {
	var ops []string

	moveMinToTop(a, &ops)

	stack.Pb(a, b)
	ops = append(ops, "pb")

	ops = append(ops, SortThree(a)...)

	stack.Pa(a, b)
	ops = append(ops, "pa")

	return ops
}

func SortFive(a, b *stack.Stack) []string {
	var ops []string

	for i := 0; i < 2; i++ {
		moveMinToTop(a, &ops)

		stack.Pb(a, b)
		ops = append(ops, "pb")
	}

	ops = append(ops, SortThree(a)...)

	stack.Pa(a, b)
	ops = append(ops, "pa")

	stack.Pa(a, b)
	ops = append(ops, "pa")

	return ops
}