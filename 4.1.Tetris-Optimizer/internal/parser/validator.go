package parser

import "fmt"

func Validate(t Tetromino) error {

	if len(t.Raw) != 4 {
		return fmt.Errorf("invalid row count")
	}

	hashCount := 0

	for _, row := range t.Raw {

		if len(row) != 4 {
			return fmt.Errorf("invalid column count")
		}

		for _, ch := range row {

			if ch != '#' && ch != '.' {
				return fmt.Errorf("invalid character")
			}

			if ch == '#' {
				hashCount++
			}
		}
	}

	if hashCount != 4 {
		return fmt.Errorf("tetromino must contain exactly 4 blocks")
	}

	return nil
}