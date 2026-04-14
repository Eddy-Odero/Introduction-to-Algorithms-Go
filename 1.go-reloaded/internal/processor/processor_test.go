package processor

import "testing"
func TestProcessLine(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		
		{"This is so exciting (up, 2)", "This is SO EXCITING"},

		{"hello world (up) (low)", "hello world"},

		{"one two three (up, 5)", "ONE TWO THREE"},

		{"FF (hex) and 1010 (bin)", "255 and 10"},

		{"ZZ (hex) 102 (bin)", "ZZ 102"},

		{"a apple a house A elephant", "an apple an house An elephant"},

		{"hello , world !", "hello, world!"},

		{"wait ... what !?", "wait... what!?"},

		{"hello ,world", "hello, world"},

		{"' hello '", "'hello'"},

		{"' hello world '", "'hello world'"},

		{"Punctuation tests are ... kinda boring ,what do you think ?", "Punctuation tests are... kinda boring, what do you think?"},

		{"it (cap) was the best (up, 3)", "It WAS THE BEST"},
	}

	for _, tt := range tests {
		got := ProcessLine(tt.input)
		if got != tt.want {
			t.Errorf("\ninput: %q\ngot:   %q\nwant:  %q\n", tt.input, got, tt.want)
		}
	}
}