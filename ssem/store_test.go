package ssem

import (
	"strconv"
	"strings"
	"testing"
)

func TestString(t *testing.T) {
	s := new(Store)
	result := s.String()
	expected := strings.Repeat("00000000000000000000000000000000\n", 32)
	if result != expected {
		t.Errorf("string repr of an empty store")
	}
}

func TestAppendBinary(t *testing.T) {
	expectedList := [...]string{
		"00000000000000000000000000000000",
		"10000000000000000000000000000000",
		"00100000000000000000000000000000",
		"00000000000000000000000000000001",
		"11111111111111111111111111111111",
		"11111111111111111111111111111110",
		"01111111111111111111111111111111",
		"11111111111111110000000000000000",
		"00000000000000001111111111111111",
		"01010101010101010101010101010101",
		"10101010101010101010101010101010"}

	for _, expected := range expectedList {
		b := strings.Builder{}

		// Create a word from the expected string (copied from ReadSnp())
		w, _ := strconv.ParseUint(Reverse(expected), 2, 33)
		word := Word(w)

		// Print it
		AppendBinary(&b, word)

		if b.String() != expected {
			t.Errorf("%s != %s", b.String(), expected)
		}
	}
}

func TestReverse(t *testing.T) {
	testMap := map[string]string{
		"":                                 "",
		"a":                                "a",
		"input":                            "tupni",
		"11111111111111110000000000000000": "00000000000000001111111111111111",
		"01010101010101010101010101010101": "10101010101010101010101010101010"}

	for input, expected := range testMap {
		if Reverse(input) != expected {
			t.Errorf("%s != %s", Reverse(input), expected)
		}
	}
}
