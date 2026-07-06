package parser

import (
	"errors"
	"strconv"
	"strings"

	"lem-in/internal/graph"
)

// Parse turns raw file lines into a validated Colony, or returns an error
// matching the spec's required "ERROR: invalid data format" wording.
func Parse(lines []string) (*graph.Colony, error) {
	colony := graph.NewColony()
	colony.RawLines = lines // keep original content for later display

	// --- Step 1: find the ant count (first non-empty, non-comment line) ---
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

	// --- Step 2: walk the rest of the lines, classifying each ---
	nextIsStart := false
	nextIsEnd := false
	// links can reference rooms defined anywhere, but every real-world lem-in
	// file lists rooms before links, so we defer link creation to a second
	// pass to avoid "unknown room" false positives caused purely by ordering.
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
			continue // ordinary comment, ignore

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

	// --- Step 3: validate start/end exist ---
	if colony.Start == nil {
		return nil, errors.New("ERROR: invalid data format, no start room found")
	}
	if colony.End == nil {
		return nil, errors.New("ERROR: invalid data format, no end room found")
	}

	// --- Step 4: now wire up links, since all rooms are known ---
	for _, pl := range pendingLinks {
		if !colony.AddLink(pl.a, pl.b) {
			return nil, errors.New("ERROR: invalid data format, link to unknown room")
		}
	}

	return colony, nil
}

// isLinkLine reports whether a line looks like "name1-name2" rather than a
// room definition. A link line has exactly one '-' and no whitespace.
func isLinkLine(line string) bool {
	if strings.ContainsAny(line, " \t") {
		return false
	}
	parts := strings.Split(line, "-")
	return len(parts) == 2 && parts[0] != "" && parts[1] != ""
}
