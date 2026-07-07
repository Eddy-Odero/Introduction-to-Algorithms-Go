package simulator

import (
	"strings"
	"testing"

	"lem-in/internal/graph"
	"lem-in/internal/solver"
)

func makePath(names ...string) []*graph.Room {
	rooms := make([]*graph.Room, len(names))
	for i, n := range names {
		rooms[i] = &graph.Room{Name: n}
	}
	return rooms
}

func TestSimulateTurnCountSinglePath(t *testing.T) {
	path := makePath("start", "a", "b", "end") // 3 edges
	assignments := solver.AssignAnts(4, [][]*graph.Room{path})
	turns := Simulate(assignments)
	if len(turns) != 6 { // 3 edges + 4 ants - 1
		t.Fatalf("expected 6 turns, got %d: %v", len(turns), turns)
	}
}

func TestSimulateTurnCountTwoPaths(t *testing.T) {
	pathA := makePath("start", "h", "n", "e", "end")
	pathB := makePath("start", "t", "E", "a", "m", "end")
	assignments := solver.AssignAnts(10, [][]*graph.Room{pathA, pathB})
	turns := Simulate(assignments)
	if len(turns) != 9 {
		t.Fatalf("expected 9 turns, got %d: %v", len(turns), turns)
	}
}

func TestSimulateEveryAntReachesEndExactlyOnce(t *testing.T) {
	path := makePath("start", "a", "end")
	assignments := solver.AssignAnts(3, [][]*graph.Room{path})
	turns := Simulate(assignments)

	endArrivals := 0
	for _, line := range turns {
		for _, tok := range strings.Fields(line) {
			if strings.HasSuffix(tok, "-end") {
				endArrivals++
			}
		}
	}
	if endArrivals != 3 {
		t.Fatalf("expected exactly 3 ants to arrive at 'end' (once each), got %d", endArrivals)
	}
}

func TestSimulateNoAntSkipsARoom(t *testing.T) {
	// Ant 1 on a 3-edge path must appear exactly once per turn until it
	// reaches "end", then never again -- it can't skip rooms or reappear.
	path := makePath("start", "a", "b", "end")
	assignments := solver.AssignAnts(1, [][]*graph.Room{path})
	turns := Simulate(assignments)

	expectedRooms := []string{"a", "b", "end"}
	if len(turns) != len(expectedRooms) {
		t.Fatalf("expected %d turns for a single ant on a 3-edge path, got %d", len(expectedRooms), len(turns))
	}
	for i, line := range turns {
		want := "L1-" + expectedRooms[i]
		if strings.TrimSpace(line) != want {
			t.Fatalf("turn %d: expected %q, got %q", i+1, want, line)
		}
	}
}
