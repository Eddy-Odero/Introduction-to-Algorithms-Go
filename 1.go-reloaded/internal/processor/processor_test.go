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
		// (cap)
{"welcome to the brooklyn bridge (cap)", "welcome to the brooklyn Bridge"},
{"this is so exciting (cap, 3)", "this Is So Exciting"},

// (low, N)
{"THIS IS SO EXCITING (low, 2)", "THIS IS so exciting"},

// article before capital vowel
{"a Elephant in the room", "an Elephant in the room"},
{"a Orange", "an Orange"},

// (cap) chained
{"hello world (cap) (up)", "hello WORLD"},

// punctuation groups at edges
{"wow !!", "wow!!"},
{"I was thinking ... You were right", "I was thinking... You were right"},
{"... hello", "...hello"},
{"hello ...", "hello..."},

// empty string
{"", ""},
	}

	for _, tt := range tests {
		got := ProcessLine(tt.input)
		if got != tt.want {
			t.Errorf("\ninput: %q\ngot:   %q\nwant:  %q\n", tt.input, got, tt.want)
		}
	}
}