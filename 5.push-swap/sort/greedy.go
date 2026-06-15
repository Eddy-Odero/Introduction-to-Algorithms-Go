package sort

import (
	"push-swap/stack"
	"push-swap/utils"
)

func GreedySort(a, b *stack.Stack) []string {
	var ops []string

	lisVals := lisValueSet(a.Data)
	size := len(a.Data)

	for i := 0; i < size; i++ {
		if lisVals[a.Data[0]] {
			stack.Ra(a)
			ops = append(ops, "ra")
		} else {
			stack.Pb(a, b)
			ops = append(ops, "pb")
		}
	}

	ops = append(ops, rotateToMin(a)...)

	for len(b.Data) > 0 {
		bestCost := int(^uint(0) >> 1)
		bestBI := 0

		for bi := range b.Data {
			aPos := insertPos(a.Data, b.Data[bi])
			c := moveCost(bi, len(b.Data), aPos, len(a.Data))
			if c < bestCost {
				bestCost = c
				bestBI = bi
			}
		}

		val := b.Data[bestBI]
		aPos := insertPos(a.Data, val)
		ops = append(ops, applyMoves(a, b, bestBI, aPos)...)

		stack.Pa(a, b)
		ops = append(ops, "pa")
	}

	ops = append(ops, rotateToMin(a)...)

	return ops
}

func applyMoves(a, b *stack.Stack, bIdx, aPos int) []string {
	var ops []string

	bSize := len(b.Data)
	aSize := len(a.Data)

	bFwd := bIdx
	bRev := bSize - bIdx
	aFwd := aPos
	aRev := aSize - aPos

	// pick cheapest combo of directions
	type combo struct {
		cost         int
		bDir, aDir   int // 1=forward, -1=reverse
	}
	imax := func(x, y int) int {
		if x > y {
			return x
		}
		return y
	}
	combos := []combo{
		{imax(bFwd, aFwd), 1, 1},
		{imax(bRev, aRev), -1, -1},
		{bFwd + aRev, 1, -1},
		{bRev + aFwd, -1, 1},
	}
	best := combos[0]
	for _, c := range combos[1:] {
		if c.cost < best.cost {
			best = c
		}
	}

	bSteps := bFwd
	if best.bDir < 0 {
		bSteps = bRev
	}
	aSteps := aFwd
	if best.aDir < 0 {
		aSteps = aRev
	}

	shared := 0
	if best.bDir == best.aDir {
		if bSteps < aSteps {
			shared = bSteps
		} else {
			shared = aSteps
		}
		if best.bDir > 0 {
			for i := 0; i < shared; i++ {
				stack.Ra(a)
				stack.Rb(b)
				ops = append(ops, "rr")
			}
		} else {
			for i := 0; i < shared; i++ {
				stack.Rra(a)
				stack.Rrb(b)
				ops = append(ops, "rrr")
			}
		}
	}

	for i := 0; i < aSteps-shared; i++ {
		if best.aDir > 0 {
			stack.Ra(a)
			ops = append(ops, "ra")
		} else {
			stack.Rra(a)
			ops = append(ops, "rra")
		}
	}

	for i := 0; i < bSteps-shared; i++ {
		if best.bDir > 0 {
			stack.Rb(b)
			ops = append(ops, "rb")
		} else {
			stack.Rrb(b)
			ops = append(ops, "rrb")
		}
	}

	return ops
}

func insertPos(a []int, val int) int {
	n := len(a)
	if n == 0 {
		return 0
	}

	maxIdx := 0
	for i, v := range a {
		if v > a[maxIdx] {
			maxIdx = i
		}
	}

	minIdx := (maxIdx + 1) % n

	if val > a[maxIdx] {
		return (maxIdx + 1) % n
	}
	if val < a[minIdx] {
		return minIdx
	}

	for i := 0; i < n; i++ {
		next := (i + 1) % n
		if a[i] < val && val <= a[next] {
			return next
		}
	}

	return minIdx
}

func moveCost(bIdx, bSize, aPos, aSize int) int {
	bFwd := bIdx
	bRev := bSize - bIdx
	aFwd := aPos
	aRev := aSize - aPos

	imax := func(x, y int) int {
		if x > y {
			return x
		}
		return y
	}
	imin := func(x, y int) int {
		if x < y {
			return x
		}
		return y
	}

	costs := [4]int{
		imax(bFwd, aFwd),
		imax(bRev, aRev),
		bFwd + aRev,
		bRev + aFwd,
	}
	best := costs[0]
	for _, c := range costs[1:] {
		best = imin(best, c)
	}
	return best
}

func rotateToMin(a *stack.Stack) []string {
	var ops []string
	if len(a.Data) <= 1 {
		return ops
	}

	minIdx := utils.IndexOf(a.Data, utils.Min(a.Data))
	n := len(a.Data)

	if minIdx <= n/2 {
		for i := 0; i < minIdx; i++ {
			stack.Ra(a)
			ops = append(ops, "ra")
		}
	} else {
		for i := 0; i < n-minIdx; i++ {
			stack.Rra(a)
			ops = append(ops, "rra")
		}
	}
	return ops
}

func lisValueSet(nums []int) map[int]bool {
	n := len(nums)
	if n == 0 {
		return nil
	}

	tails := []int{}
	prev := make([]int, n)
	tailIdx := make([]int, n)
	for i := range prev {
		prev[i] = -1
	}

	for i := 0; i < n; i++ {
		lo, hi := 0, len(tails)
		for lo < hi {
			mid := (lo + hi) / 2
			if nums[tails[mid]] < nums[i] {
				lo = mid + 1
			} else {
				hi = mid
			}
		}
		tailIdx[i] = lo
		if lo == len(tails) {
			tails = append(tails, i)
		} else {
			tails[lo] = i
		}
		if lo > 0 {
			prev[i] = tails[lo-1]
		}
	}

	set := make(map[int]bool)
	idx := tails[len(tails)-1]
	for idx != -1 {
		set[nums[idx]] = true
		idx = prev[idx]
	}
	return set
}