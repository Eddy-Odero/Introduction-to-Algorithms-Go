package processor

import "testing"

func TestProcessLine(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello , world !", "hello, world!"},
		{"hello world (up, 2)", "HELLO WORLD"},
	}

	for _, tt := range tests {
		got := ProcessLine(tt.input)

		if got != tt.want {
			t.Errorf("input: %q | got: %q | want: %q",
				tt.input, got, tt.want)
		}
	}
}