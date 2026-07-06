package solver

import "lem-in/internal/graph"

type PathAssignment struct {
	Path   []*graph.Room
	AntIDs []int
}

func AssignAnts(numAnts int, paths [][]*graph.Room) []*PathAssignment {
	assignments := make([]*PathAssignment, len(paths))
	load := make([]int, len(paths)) // projected finish turn if we add one more ant

	for i, p := range paths {
		edges := len(p) - 1
		load[i] = edges // turns if this path carries 1 ant so far... see below
		assignments[i] = &PathAssignment{Path: p}
	}

	for ant := 1; ant <= numAnts; ant++ {
		best := 0
		for i := 1; i < len(load); i++ {
			if load[i] < load[best] {
				best = i
			}
		}
		assignments[best].AntIDs = append(assignments[best].AntIDs, ant)
		load[best]++ // one more ant on this path pushes its finish turn out by 1
	}

	return assignments
}