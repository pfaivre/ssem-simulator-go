package main

import (
	"fmt"
	"strings"
)

// Number of words in the store
const word_count = 32

// Position of the address an instruction
const address_start = 0

// Size of addresses in bits
const address_length = 5

// Position of the operation code in the word
const opcode_start = 13

// Operation code size in bits
const opcode_length = 3

// A word is representaed by a 32 bits integer
type Word int32

// Contains the machine's memory
type Store [word_count]Word

func (s Store) String() string {
	builder := strings.Builder{}
	for _, w := range s {
		builder.WriteString(fmt.Sprintf("%032b\n", w))
	}
	return builder.String()
}

// A machine has a store, a counter increment (ci) and an accumulator register (a)
type Ssem struct {
	store Store
	ci    Word
	a     Word
}

func (s Ssem) String() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%032b CI\n", s.ci))
	builder.WriteString(fmt.Sprintf("%032b A\n\n", s.a))
	builder.WriteString(fmt.Sprint(s.store))
	return builder.String()
}

func main() {
	fmt.Println("SSEM Simulator, 2025 Pierre Faivre")
	fmt.Println("(Work in progress, not functional yet)")
	fmt.Println()
	machine := Ssem{}
	fmt.Println(machine)
}
