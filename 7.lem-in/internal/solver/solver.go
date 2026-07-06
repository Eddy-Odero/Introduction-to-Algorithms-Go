package solver

import "lem-in/internal/graph"

// PathAssignment holds one discovered path plus the specific ant numbers
// (1-indexed, matching the output spec) that will travel on it, in the
// order they enter Start.
type PathAssignment struct {
	Path   []*graph.Room
	AntIDs []int
}

// AssignAnts distributes numAnts across the given disjoint paths to
// minimize the total number of turns needed for every ant to reach End.
//
// Turns for a path carrying k ants = (edges in path) + k - 1.
// Since every ant is "equal work", the optimal distribution is greedy:
// send each next ant down whichever path currently has the lowest
// projected finish time. This is a straightforward exchange-argument
// optimum for identical-size jobs across parallel "machines" with
// different fixed head-starts (the path lengths).
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