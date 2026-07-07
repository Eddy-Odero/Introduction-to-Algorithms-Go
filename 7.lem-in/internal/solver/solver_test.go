package solver

import (
	"testing"

	"lem-in/internal/graph"
)

func makePath(names ...string) []*graph.Room {
	rooms := make([]*graph.Room, len(names))
	for i, n := range names {
		rooms[i] = &graph.Room{Name: n}
	}
	return rooms
}

func TestAssignAntsSinglePathGetsAllAnts(t *testing.T) {
	path := makePath("start", "a", "b", "end") // 3 edges
	assignments := AssignAnts(4, [][]*graph.Room{path})
	if len(assignments[0].AntIDs) != 4 {
		t.Fatalf("expected all 4 ants on the only available path, got %d", len(assignments[0].AntIDs))
	}
}

func TestAssignAntsBalancesTwoPathsOptimally(t *testing.T) {
	pathA := makePath("start", "h", "n", "e", "end")     // 4 edges
	pathB := makePath("start", "t", "E", "a", "m", "end") // 5 edges
	assignments := AssignAnts(10, [][]*graph.Room{pathA, pathB})

	kA, kB := len(assignments[0].AntIDs), len(assignments[1].AntIDs)
	if kA+kB != 10 {
		t.Fatalf("expected all 10 ants assigned, got %d", kA+kB)
	}
	turnsA, turnsB := 4+kA-1, 5+kB-1
	maxTurns := turnsA
	if turnsB > maxTurns {
		maxTurns = turnsB
	}
	// Known-correct answer for this classic split: 9 turns.
	if maxTurns != 9 {
		t.Fatalf("expected optimal max turns of 9, got %d (kA=%d, kB=%d)", maxTurns, kA, kB)
	}
}

func TestAssignAntsEveryAntUsedExactlyOnce(t *testing.T) {
	pathA := makePath("start", "a", "end")
	pathB := makePath("start", "b", "c", "end")
	assignments := AssignAnts(6, [][]*graph.Room{pathA, pathB})

	seen := map[int]bool{}
	for _, a := range assignments {
		for _, id := range a.AntIDs {
			if seen[id] {
				t.Fatalf("ant %d assigned to more than one path", id)
			}
			seen[id] = true
		}
	}
	for i := 1; i <= 6; i++ {
		if !seen[i] {
			t.Fatalf("ant %d was never assigned to any path", i)
		}
	}
}
