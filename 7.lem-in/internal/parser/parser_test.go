package parser

import (
	"strings"
	"testing"
)

func mustParse(t *testing.T, input string) {
	t.Helper()
	if _, err := Parse(strings.Split(input, "\n")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func mustFail(t *testing.T, input string) {
	t.Helper()
	if _, err := Parse(strings.Split(input, "\n")); err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestParseValidColony(t *testing.T) {
	input := `4
##start
0 0 3
2 2 5
3 4 0
##end
1 8 3
0-2
2-3
3-1`
	colony, err := Parse(strings.Split(input, "\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if colony.NumAnts != 4 {
		t.Errorf("expected 4 ants, got %d", colony.NumAnts)
	}
	if colony.Start.Name != "0" {
		t.Errorf("expected start room '0', got %q", colony.Start.Name)
	}
	if colony.End.Name != "1" {
		t.Errorf("expected end room '1', got %q", colony.End.Name)
	}
	if len(colony.Rooms) != 4 {
		t.Errorf("expected 4 rooms, got %d", len(colony.Rooms))
	}
}

func TestParseCommentsAreIgnored(t *testing.T) {
	input := `2
#comment before start
##start
0 0 0
#comment between rooms
##end
1 1 1
0-1`
	colony, err := Parse(strings.Split(input, "\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if colony.NumAnts != 2 {
		t.Errorf("expected 2 ants, got %d", colony.NumAnts)
	}
}

func TestParseInvalidAntCountNonNumeric(t *testing.T) {
	mustFail(t, "abc\n##start\n0 0 0\n##end\n1 1 1\n0-1")
}

func TestParseInvalidAntCountZero(t *testing.T) {
	mustFail(t, "0\n##start\n0 0 0\n##end\n1 1 1\n0-1")
}

func TestParseInvalidAntCountNegative(t *testing.T) {
	mustFail(t, "-3\n##start\n0 0 0\n##end\n1 1 1\n0-1")
}

func TestParseDuplicateRoomName(t *testing.T) {
	mustFail(t, "2\n##start\n0 0 0\n##end\n1 1 1\n0 2 2\n0-1")
}

func TestParseLinkToUnknownRoom(t *testing.T) {
	mustFail(t, "2\n##start\n0 0 0\n##end\n1 1 1\n0-9")
}

func TestParseMissingStartRoom(t *testing.T) {
	mustFail(t, "2\n0 0 0\n##end\n1 1 1\n0-1")
}

func TestParseMissingEndRoom(t *testing.T) {
	mustFail(t, "2\n##start\n0 0 0\n1 1 1\n0-1")
}

func TestParseDuplicateStartMarker(t *testing.T) {
	mustFail(t, "2\n##start\n0 0 0\n##start\n1 1 1\n##end\n2 2 2\n0-1\n1-2")
}

func TestParseMalformedRoomLine(t *testing.T) {
	mustFail(t, "2\n##start\n0 0\n##end\n1 1 1\n0-1") // missing y coordinate
}

func TestParseInvalidRoomCoordinates(t *testing.T) {
	mustFail(t, "2\n##start\n0 x y\n##end\n1 1 1\n0-1")
}

func TestParseLinksBeforeRoomsStillWork(t *testing.T) {
	// Links reference rooms that are only defined later in the file --
	// should still succeed since link wiring is deferred to a second pass.
	mustParse(t, "2\n0-1\n##start\n0 0 0\n##end\n1 1 1")
}
