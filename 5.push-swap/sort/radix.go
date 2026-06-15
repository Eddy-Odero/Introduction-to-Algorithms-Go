package sort

import (
	"push-swap/stack"
	"push-swap/utils"
)
func RadixSort(a, b *stack.Stack) []string {
	var ops []string

	maxBits := utils.MaxBits(a.Data)

	for bit := 0; bit < maxBits; bit++ {

		size := len(a.Data)

		for i := 0; i < size; i++ {

			if len(a.Data) == 0 {
				break
			}

			num := a.Data[0]

			if (num>>bit)&1 == 1 {
				stack.Ra(a)
				ops = append(ops, "ra")
			} else {
				stack.Pb(a, b)
				ops = append(ops, "pb")
			}
		}

		for len(b.Data) > 0 {
			stack.Pa(a, b)
			ops = append(ops, "pa")
		}
	}

	return ops
}