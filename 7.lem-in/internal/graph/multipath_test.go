package graph

import "testing"

func TestFindAllPathsFindsDisjointRoutes(t *testing.T) {
	// Two equal-length disjoint routes from start to end:
	// start-a-b-end and start-c-d-end
	c := NewColony()
	for _, n := range []string{"start", "a", "b", "end", "c", "d"} {
		c.AddRoom(n, 0, 0)
	}
	c.AddLink("start", "a")
	c.AddLink("a", "b")
	c.AddLink("b", "end")
	c.AddLink("start", "c")
	c.AddLink("c", "d")
	c.AddLink("d", "end")
	c.Start = c.Rooms["start"]
	c.End = c.Rooms["end"]

	paths := FindAllPaths(c)
	if len(paths) != 2 {
		t.Fatalf("expected 2 disjoint paths, got %d", len(paths))
	}

	// no intermediate room (excluding start/end) should be used by more than one path
	seen := map[string]int{}
	for _, p := range paths {
		for _, r := range p {
			if r.Name == "start" || r.Name == "end" {
				continue
			}
			seen[r.Name]++
		}
	}
	for name, count := range seen {
		if count > 1 {
			t.Fatalf("room %q reused across paths %d times; paths must be disjoint", name, count)
		}
	}
}

func TestFindAllPathsSingleRouteOnly(t *testing.T) {
	c := buildChain() // simple 0-1-2-3 chain, no alternate route
	paths := FindAllPaths(c)
	if len(paths) != 1 {
		t.Fatalf("expected exactly 1 path for a simple chain, got %d", len(paths))
	}
}

func TestFindAllPathsNoPath(t *testing.T) {
	c := NewColony()
	c.AddRoom("start", 0, 0)
	c.AddRoom("end", 1, 1)
	c.Start = c.Rooms["start"]
	c.End = c.Rooms["end"]
	paths := FindAllPaths(c)
	if len(paths) != 0 {
		t.Fatalf("expected 0 paths for disconnected start/end, got %d", len(paths))
	}
}
