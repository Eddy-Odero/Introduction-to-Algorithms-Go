package parse

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseArgs(args []string) ([]int, error) {
	var values []int

	seen := map[int]bool{}

	for _, arg := range args {
		fields := strings.Fields(arg)

		for _, field := range fields {
			n, err := strconv.Atoi(field)

			if err != nil {
				return nil, fmt.Errorf("invalid integer")
			}

			if seen[n] {
				return nil, fmt.Errorf("duplicate value")
			}

			seen[n] = true

			values = append(values, n)
		}
	}

	return values, nil
}