package graph

import "testing"

func TestAddRoomDuplicate(t *testing.T) {
	c := NewColony()
	if _, ok := c.AddRoom("1", 0, 0); !ok {
		t.Fatal("expected first AddRoom to succeed")
	}
	if _, ok := c.AddRoom("1", 1, 1); ok {
		t.Fatal("expected duplicate room name to fail")
	}
}

func TestAddLinkUnknownRoom(t *testing.T) {
	c := NewColony()
	c.AddRoom("1", 0, 0)
	if c.AddLink("1", "2") {
		t.Fatal("expected link to an undefined room to fail")
	}
}

func TestAddLinkBidirectional(t *testing.T) {
	c := NewColony()
	c.AddRoom("1", 0, 0)
	c.AddRoom("2", 1, 1)
	if !c.AddLink("1", "2") {
		t.Fatal("expected link between two existing rooms to succeed")
	}
	r1, r2 := c.Rooms["1"], c.Rooms["2"]
	if len(r1.Links) != 1 || r1.Links[0] != r2 {
		t.Fatal("room 1 should have room 2 as a neighbor")
	}
	if len(r2.Links) != 1 || r2.Links[0] != r1 {
		t.Fatal("tunnels are bidirectional: room 2 should link back to room 1")
	}
}
