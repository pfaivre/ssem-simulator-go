package ssem

import (
	"fmt"
	"strings"
)

// A word is represented by a 32 bits integer
type Word int32

// Contains the machine's memory
type Store [WORD_COUNT]Word

// Reads the address pointed at the given address and parses its given operation code and data
func (s *Store) DecodeInstruction(address Word) (Opcode, Word) {
	word := s[address]

	// Objective: extract opcode and data from word
	// word: 0b00000000000000000100000000011000
	//                         ===        =====
	//        Operation code ---'           |
	//                  Data ---------------'

	// Step 1: extract instruction data (5 bits)
	// word: 0b00000000000000000100000000011000
	// mask: 0b00000000000000000000000000011111
	//       ----------------------------------
	//    &: 0b00000000000000000000000000011000
	data := SSEM_DATA_MASK & word

	// Step 2: Shift bits to put the opcode on the edge
	// word: 0b00000000000000000100000000011000
	//       ----------------------------------
	// >>13: 0b00000000000000000000000000000010

	// Step 3: Extract opcode (3 bits)
	// word: 0b00000000000000000000000000000010
	// mask: 0b00000000000000000000000000000111
	//       ----------------------------------
	//    &: 0b00000000000000000000000000000010
	opcode := (SSEM_OPCODE_MASK & (word >> OPCODE_START))

	return Opcode(opcode), Word(data)
}

func (s Store) String() string {
	builder := strings.Builder{}
	builder.Grow(4096) // takes approximately 2304 bytes to print the store
	for i, w := range s {
		AppendBinary(&builder, w)
		o, d := s.DecodeInstruction(Word(i))
		builder.WriteString(fmt.Sprintf(" %02d: %11d  %s %02d   \n", i, w, o, d))
	}

	return builder.String()
}

var BinaryDigitReplacer = strings.NewReplacer()

func AppendBinary(b *strings.Builder, w Word) {
	b.WriteString("\033[01;32m")
	b.WriteString(
		BinaryDigitReplacer.Replace(
			Reverse(fmt.Sprintf("%032b", uint32(w)))))
	b.WriteString("\033[0m")
}

// Returns a new string with each rune in reverse order
func Reverse(s string) string {
	// TODO: consider swapping slices of bytes instead for more efficiency, if possible
	// func Reverse(s []byte)
	var b strings.Builder
	b.Grow(len(s)) // Allocate the definitive size from the beginning to avoid multiple allocations
	runes := []rune(s)
	for i := len(runes) - 1; i >= 0; i-- {
		b.WriteRune(runes[i])
	}
	return b.String()
}
