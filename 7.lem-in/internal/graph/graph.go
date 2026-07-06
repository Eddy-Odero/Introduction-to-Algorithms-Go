package graph

// Room represents a single room in the ant colony.
// Links holds pointers to directly-connected rooms (shared instances,
// not copies) so traversal later is a simple `for _, n := range room.Links`.
type Room struct {
	Name  string
	X, Y  int
	Links []*Room
}

// Colony represents the whole parsed, validated ant farm:
// every room, which one is start/end, and how many ants need to cross.
type Colony struct {
	NumAnts  int
	Rooms    map[string]*Room // name -> Room, for O(1) lookup while parsing links
	Start    *Room
	End      *Room
	RawLines []string // original file lines, kept for the required output display
}

// NewColony returns an empty Colony ready to be filled in by the parser.
func NewColony() *Colony {
	return &Colony{
		Rooms: make(map[string]*Room),
	}
}

// AddRoom creates a Room, stores it in the colony, and returns it.
// Returns false if a room with that name already exists (duplicate rooms
// are invalid input per spec).
func (c *Colony) AddRoom(name string, x, y int) (*Room, bool) {
	if _, exists := c.Rooms[name]; exists {
		return nil, false
	}
	r := &Room{Name: name, X: x, Y: y}
	c.Rooms[name] = r
	return r, true
}

// AddLink connects two existing rooms bidirectionally (tunnels work both ways).
// Returns false if either room name is unknown.
func (c *Colony) AddLink(name1, name2 string) bool {
	r1, ok1 := c.Rooms[name1]
	r2, ok2 := c.Rooms[name2]
	if !ok1 || !ok2 {
		return false
	}
	r1.Links = append(r1.Links, r2)
	r2.Links = append(r2.Links, r1)
	return true
}
