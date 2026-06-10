package utils

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