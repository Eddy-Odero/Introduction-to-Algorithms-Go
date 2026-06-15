package utils

import "sort"

func IsSorted(nums []int) bool {
	for i := 0; i < len(nums)-1; i++ {
		if nums[i] > nums[i+1] {
			return false
		}
	}
	return true
}

func Min(nums []int) int {
	min := nums[0]

	for _, n := range nums {
		if n < min {
			min = n
		}
	}

	return min
}

func IndexOf(nums []int, value int) int {
	for i, n := range nums {
		if n == value {
			return i
		}
	}
	return -1
}

func Normalize(nums []int) []int {
	sorted := make([]int, len(nums))
	copy(sorted, nums)

	sort.Ints(sorted)

	ranks := make(map[int]int)

	for i, v := range sorted {
		ranks[v] = i
	}

	result := make([]int, len(nums))

	for i, v := range nums {
		result[i] = ranks[v]
	}

	return result
}
func MaxBits(nums []int) int {
	max := nums[0]

	for _, n := range nums {
		if n > max {
			max = n
		}
	}

	bits := 0

	for max > 0 {
		bits++
		max >>= 1
	}

	return bits
}