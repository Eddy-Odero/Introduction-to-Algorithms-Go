package graph

type Room struct {
	Name  string
	X, Y  int
	Links []*Room
}

type Colony struct {
	NumAnts  int
	Rooms    map[string]*Room 
	Start    *Room
	End      *Room
	RawLines []string 
}

func NewColony() *Colony {
	return &Colony{
		Rooms: make(map[string]*Room),
	}
}

func (c *Colony) AddRoom(name string, x, y int) (*Room, bool) {
	if _, exists := c.Rooms[name]; exists {
		return nil, false
	}
	r := &Room{Name: name, X: x, Y: y}
	c.Rooms[name] = r
	return r, true
}

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
