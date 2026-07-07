package graph

import "testing"

// buildChain creates a simple 4-room straight-line colony: 0-1-2-3
func buildChain() *Colony {
	c := NewColony()
	for i, n := range []string{"0", "1", "2", "3"} {
		c.AddRoom(n, i, 0)
	}
	c.AddLink("0", "1")
	c.AddLink("1", "2")
	c.AddLink("2", "3")
	c.Start = c.Rooms["0"]
	c.End = c.Rooms["3"]
	return c
}

func namesOf(rooms []*Room) []string {
	out := make([]string, len(rooms))
	for i, r := range rooms {
		out[i] = r.Name
	}
	return out
}

func TestFindShortestPathChain(t *testing.T) {
	c := buildChain()
	path := FindShortestPath(c)
	if path == nil {
		t.Fatal("expected a path to be found")
	}
	want := []string{"0", "1", "2", "3"}
	got := namesOf(path)
	if len(got) != len(want) {
		t.Fatalf("expected path %v, got %v", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected path %v, got %v", want, got)
		}
	}
}

func TestFindShortestPathNoConnection(t *testing.T) {
	c := NewColony()
	c.AddRoom("a", 0, 0)
	c.AddRoom("b", 1, 1)
	c.Start = c.Rooms["a"]
	c.End = c.Rooms["b"]
	if path := FindShortestPath(c); path != nil {
		t.Fatalf("expected nil path for disconnected rooms, got %v", namesOf(path))
	}
}

func TestFindShortestPathPrefersShorterRoute(t *testing.T) {
	// start-a-end is 2 edges; start-b-c-end is 3 edges. BFS must pick the former.
	c := NewColony()
	for _, n := range []string{"start", "a", "end", "b", "c"} {
		c.AddRoom(n, 0, 0)
	}
	c.AddLink("start", "a")
	c.AddLink("a", "end")
	c.AddLink("start", "b")
	c.AddLink("b", "c")
	c.AddLink("c", "end")
	c.Start = c.Rooms["start"]
	c.End = c.Rooms["end"]

	path := FindShortestPath(c)
	if len(path) != 3 {
		t.Fatalf("expected shortest 3-room path (start-a-end), got %v", namesOf(path))
	}
}
