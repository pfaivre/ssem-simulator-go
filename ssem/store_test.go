package ssem

import (
	"strings"
	"testing"
)

const ANSI_COLOR = "\033[01;32m"
const ANSI_RESET = "\033[0m"

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

func TestAppendBinary(t *testing.T) {
	// The SSEM stores the least significant bit first, so the binary rendering
	// is the reverse of the usual MSB-first notation.
	tests := []struct {
		word Word
		want string
	}{
		{0, "00000000000000000000000000000000"},
		{1, "10000000000000000000000000000000"},
		{2, "01000000000000000000000000000000"},
		{3, "11000000000000000000000000000000"},
		{-1, "11111111111111111111111111111111"},
		{-2, "01111111111111111111111111111111"},
		{57344, "00000000000001110000000000000000"},
	}

	for _, tt := range tests {
		b := strings.Builder{}
		AppendBinary(&b, tt.word)

		expected := ANSI_COLOR + tt.want + ANSI_RESET
		if b.String() != expected {
			t.Errorf("AppendBinary(%d) = %q, want %q", tt.word, b.String(), expected)
		}
	}
}

func TestStoreWrite(t *testing.T) {
	store := Store{}
	store[0] = Word(STP) << OPCODE_START    // STP 00
	store[5] = Word(LDN)<<OPCODE_START | 31 // LDN 31

	b := strings.Builder{}
	store.Write(&b, 5)
	out := b.String()

	lines := strings.Split(out, "\n")
	if len(lines) != WORD_COUNT+1 {
		t.Errorf("Write produced %d lines, want %d", len(lines)-1, WORD_COUNT)
	}

	if !strings.HasPrefix(lines[0], " ") {
		t.Errorf("line 0 should start with a space, got %q", lines[0])
	}
	if !strings.Contains(lines[0], " 00:       57344  STP 00   ") {
		t.Errorf("line 0 should contain the decoded STP instruction, got %q", lines[0])
	}

	if !strings.HasPrefix(lines[5], ">") {
		t.Errorf("line 5 should start with '>' as it is pointed at by CI, got %q", lines[5])
	}
	if !strings.Contains(lines[5], " 05:       16415  LDN 31   ") {
		t.Errorf("line 5 should contain the decoded LDN instruction, got %q", lines[5])
	}

	if strings.Count(out, ">") != 1 {
		t.Errorf("Write should mark exactly one line with '>', got %d", strings.Count(out, ">"))
	}
}
