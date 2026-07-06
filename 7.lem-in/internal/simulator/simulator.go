package simulator

import (
	"fmt"
	"strings"

	"lem-in/internal/solver"
)

func Simulate(assignments []*solver.PathAssignment) []string {
	type activeAnt struct {
		antID   int
		pathIdx int
		roomIdx int // index into assignments[pathIdx].Path; 0 = Start
	}

	var active []*activeAnt
	// nextToEnter[pathIdx] = index into AntIDs of the next ant on that path
	// that hasn't entered the path yet.
	nextToEnter := make([]int, len(assignments))

	var turns []string

	for {
		var moves []string

		// 1. Advance ants already in transit.
		stillActive := active[:0]
		for _, a := range active {
			path := assignments[a.pathIdx].Path
			if a.roomIdx < len(path)-1 {
				a.roomIdx++
				moves = append(moves, fmt.Sprintf("L%d-%s", a.antID, path[a.roomIdx].Name))
			}
			if a.roomIdx < len(path)-1 {
				stillActive = append(stillActive, a)
			}
			// if it just reached End (last index), it's done, drop it
		}
		active = stillActive

		for pIdx, assignment := range assignments {
			if nextToEnter[pIdx] >= len(assignment.AntIDs) {
				continue // no more ants queued for this path
			}
			antID := assignment.AntIDs[nextToEnter[pIdx]]
			nextToEnter[pIdx]++
			newAnt := &activeAnt{antID: antID, pathIdx: pIdx, roomIdx: 1}
			if len(assignment.Path) > 1 {
				moves = append(moves, fmt.Sprintf("L%d-%s", antID, assignment.Path[1].Name))
			}
			if newAnt.roomIdx < len(assignment.Path)-1 {
				active = append(active, newAnt)
			}
		}

		if len(moves) == 0 {
			break // nothing moved and nothing left to enter: done
		}
		turns = append(turns, strings.Join(moves, " "))
	}

	return turns
}