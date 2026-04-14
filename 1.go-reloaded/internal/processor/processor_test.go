package processor

import "testing"

func TestProcessLine(t *testing.T) {
	got := ProcessLine("hello , world !")
	want := "hello, world!"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
