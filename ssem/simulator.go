package ssem

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Number of words in the store
const WORD_COUNT = 32

// Position of the address an instruction
const ADDRESS_START = 0

// Size of addresses in bits
const ADDRESS_LENGTH = 5

// Position of the operation code in the word
const OPCODE_START = 13

// Operation code size in bits
const OPCODE_LENGTH = 3

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

// Reinitializes the memory and loads the content of the given assembly file
func (s *Ssem) ReadAsm(file_name string) error {
	store := Store{}

	readFile, err := os.Open(file_name)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(readFile)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Split(line, ";")[0]
		line = strings.Trim(line, " ")

		if len(line) == 0 {
			continue
		}

		l := strings.Split(line, " ")

		index, err := strconv.Atoi(l[0])
		if err != nil {
			return err
		}
		if index < 0 || index > WORD_COUNT-1 {
			return fmt.Errorf("Index out of bound, expected from 0 to %d, got %d", WORD_COUNT-1, index)
		}

		mnemonic, err := FromString(l[1])
		if err != nil {
			return err
		}
		if mnemonic != NUM {
			store[index] = Word(mnemonic) << OPCODE_START
		}

		data := 0
		if len(l) > 2 {
			data, err = strconv.Atoi(l[2])
			if err != nil {
				return err
			}
		}
		store[index] |= Word(data) << ADDRESS_START
	}

	err = readFile.Close()
	if err != nil {
		return err
	}

	s.store = store
	s.ci = 0
	s.a = 0

	return nil
}
