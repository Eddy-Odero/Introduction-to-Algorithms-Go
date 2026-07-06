package parser

import (
	"errors"
	"strconv"
	"strings"

	"lem-in/internal/graph"
)

func Parse(lines []string) (*graph.Colony, error) {
	colony := graph.NewColony()
	colony.RawLines = lines 

	idx := 0
	for idx < len(lines) && (strings.TrimSpace(lines[idx]) == "" || strings.HasPrefix(lines[idx], "#")) {
		idx++
	}
	if idx >= len(lines) {
		return nil, errors.New("ERROR: invalid data format, empty file")
	}
	numAnts, err := strconv.Atoi(strings.TrimSpace(lines[idx]))
	if err != nil || numAnts <= 0 {
		return nil, errors.New("ERROR: invalid data format, invalid number of Ants")
	}
	colony.NumAnts = numAnts
	idx++

	nextIsStart := false
	nextIsEnd := false

	type pendingLink struct{ a, b string }
	var pendingLinks []pendingLink

	for ; idx < len(lines); idx++ {
		line := strings.TrimSpace(lines[idx])
		if line == "" {
			continue
		}

		switch {
		case line == "##start":
			if nextIsStart || colony.Start != nil {
				return nil, errors.New("ERROR: invalid data format, multiple start rooms")
			}
			nextIsStart = true
			continue
		case line == "##end":
			if nextIsEnd || colony.End != nil {
				return nil, errors.New("ERROR: invalid data format, multiple end rooms")
			}
			nextIsEnd = true
			continue
		case strings.HasPrefix(line, "#"):
			continue 

		case isLinkLine(line):
			parts := strings.Split(line, "-")
			pendingLinks = append(pendingLinks, pendingLink{parts[0], parts[1]})

		default: // room line: "name x y"
			fields := strings.Fields(line)
			if len(fields) != 3 {
				return nil, errors.New("ERROR: invalid data format, malformed room")
			}
			name := fields[0]
			x, errX := strconv.Atoi(fields[1])
			y, errY := strconv.Atoi(fields[2])
			if errX != nil || errY != nil {
				return nil, errors.New("ERROR: invalid data format, invalid room coordinates")
			}
			room, ok := colony.AddRoom(name, x, y)
			if !ok {
				return nil, errors.New("ERROR: invalid data format, duplicate room")
			}
			if nextIsStart {
				colony.Start = room
				nextIsStart = false
			}
			if nextIsEnd {
				colony.End = room
				nextIsEnd = false
			}
		}
	}

	if colony.Start == nil {
		return nil, errors.New("ERROR: invalid data format, no start room found")
	}
	if colony.End == nil {
		return nil, errors.New("ERROR: invalid data format, no end room found")
	}

	for _, pl := range pendingLinks {
		if !colony.AddLink(pl.a, pl.b) {
			return nil, errors.New("ERROR: invalid data format, link to unknown room")
		}
	}

	return colony, nil
}

func isLinkLine(line string) bool {
	if strings.ContainsAny(line, " \t") {
		return false
	}
	parts := strings.Split(line, "-")
	return len(parts) == 2 && parts[0] != "" && parts[1] != ""
}
