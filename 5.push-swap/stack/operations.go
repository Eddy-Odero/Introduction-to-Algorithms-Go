package stack


// SWAP


func Sa(a *Stack) {
	if len(a.Data) < 2 {
		return
	}
	a.Data[0], a.Data[1] = a.Data[1], a.Data[0]
}

func Sb(b *Stack) {
	if len(b.Data) < 2 {
		return
	}
	b.Data[0], b.Data[1] = b.Data[1], b.Data[0]
}

func Ss(a, b *Stack) {
	Sa(a)
	Sb(b)
}

// PUSH


func Pa(a, b *Stack) {
	if len(b.Data) == 0 {
		return
	}

	value := b.Data[0]
	b.Data = b.Data[1:]
	a.Data = append([]int{value}, a.Data...)
}

func Pb(a, b *Stack) {
	if len(a.Data) == 0 {
		return
	}

	value := a.Data[0]
	a.Data = a.Data[1:]
	b.Data = append([]int{value}, b.Data...)
}


// ROTATE


func Ra(a *Stack) {
	if len(a.Data) < 2 {
		return
	}

	first := a.Data[0]
	a.Data = append(a.Data[1:], first)
}

func Rb(b *Stack) {
	if len(b.Data) < 2 {
		return
	}

	first := b.Data[0]
	b.Data = append(b.Data[1:], first)
}

func Rr(a, b *Stack) {
	Ra(a)
	Rb(b)
}

// REVERSE ROTATE

func Rra(a *Stack) {
	if len(a.Data) < 2 {
		return
	}

	last := a.Data[len(a.Data)-1]
	a.Data = append([]int{last}, a.Data[:len(a.Data)-1]...)
}

func Rrb(b *Stack) {
	if len(b.Data) < 2 {
		return
	}

	last := b.Data[len(b.Data)-1]
	b.Data = append([]int{last}, b.Data[:len(b.Data)-1]...)
}

func Rrr(a, b *Stack) {
	Rra(a)
	Rrb(b)
}