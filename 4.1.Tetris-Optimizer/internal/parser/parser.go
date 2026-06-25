package parser

import (
	"strings"
	"fmt"
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

	fmt.Println("Blocks:", len(blocks))

	var pieces []Tetromino

	for _, block := range blocks {
		rows := strings.Split(block, "\n")

		pieces = append(pieces, Tetromino{
			Raw: rows,
		})
	}

	return pieces, nil
}