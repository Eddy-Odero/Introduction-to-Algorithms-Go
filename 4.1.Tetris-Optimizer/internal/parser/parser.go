package parser

import (
	"strings"
)

type Tetromino struct {
	Raw []string
}

func Parse(content string) ([]Tetromino, error) {
	content = strings.ReplaceAll(content, "\r\n", "\n")

	blocks := strings.Split(
		strings.TrimSpace(content),
		"\n\n",
	)

	

	var pieces []Tetromino

	for _, block := range blocks {

	rows := strings.Split(block, "\n")

	t := Tetromino{
		Raw: rows,
	}

	if err := Validate(t); err != nil {
		return nil, err
	}

	pieces = append(pieces, t)
}

	return pieces, nil
}