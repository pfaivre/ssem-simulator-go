package ssem

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strconv"
	"strings"
	"time"
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

const SSEM_OPCODE_MASK = 0b00000000000000000000000000000111

const SSEM_DATA_MASK = 0b00000000000000000000000000011111

// A machine has a store, a counter increment (ci) and an accumulator register (a)
type Ssem struct {
	store    Store
	ci       Word
	a        Word
	StopFlag bool
}

// Type that has a state, loaded instructions and can be executed to modify that state
type RunnableMachine interface {
	InstructionCycle() error
	DecodeInstruction(ci Word) (Opcode, Word, error)
	Run(max_cycles uint) (uint, error)
}

func (s Ssem) String() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%032b CI\n", s.ci))
	builder.WriteString(fmt.Sprintf("%032b A\n\n", s.a))
	builder.WriteString(fmt.Sprint(s.store))
	return builder.String()
}

// Reinitializes the memory and loads the content of the given assembly file
//
// ASM files have a format like this:
//
//	00 NUM 1    ;Incremental Value
//	01 LDN 31   ;Load negative of counter
//	02 SUB 0    ;"Increment" our counter
//	...
func (s *Ssem) ReadAsm(file_name string) error {
	store := Store{}

	readFile, err := os.Open(file_name)
	if err != nil {
		return err
	}
	defer readFile.Close()

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
			return fmt.Errorf("index out of bound, expected from 0 to %d, got %d", WORD_COUNT-1, index)
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

	s.store = store
	s.ci = 0
	s.a = 0
	s.StopFlag = false

	return readFile.Close()
}

// Reinitializes the memory and loads the content of the given snp file.
//
// SNP files have a format like this:
//
//	; Comment
//	00: 00000110101001000100000100000100
//	01: 10011011111100100010000010001000
//	02: 10000010000101101000100001010000
//	...
func (s *Ssem) ReadSnp(file_name string) error {
	store := Store{}

	readFile, err := os.Open(file_name)
	if err != nil {
		return err
	}
	defer readFile.Close()

	scanner := bufio.NewScanner(readFile)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Split(line, ";")[0]
		line = strings.Trim(line, " ")

		if len(line) == 0 {
			continue
		}

		l := strings.Split(line, ":")

		// Extracting index
		index, err := strconv.Atoi(l[0])
		if err != nil {
			return err
		}
		if index < 0 || index > WORD_COUNT-1 {
			return fmt.Errorf("index out of bound, expected from 0 to %d, got %d", WORD_COUNT-1, index)
		}

		// Extracting word
		if len(l) <= 1 {
			return fmt.Errorf("missing word for index %d", index)
		}
		w, err := strconv.ParseUint(strings.Trim(l[1], " "), 2, 33)
		if err != nil {
			return err
		}
		// Need to reverse the bit order as the SSEM stores number in reverse
		store[index] = Word(bits.Reverse32(uint32(w)))
	}

	s.store = store
	s.ci = 0
	s.a = 0
	s.StopFlag = false

	return readFile.Close()
}

// Performs one cycle
func (s *Ssem) InstructionCycle() error {
	// Fetch
	s.ci += 1
	s.ci %= Word(len(s.store)) // CI loops back to the begining when it exceeds the store boundaries

	// Decode
	opcode, data := s.DecodeInstruction(s.ci)

	// Execute
	err := s.Execute(opcode, data)
	if err != nil {
		return err
	}

	return nil
}

// Reads the address pointed at the given address and parses its given operation code and data
func (s *Ssem) DecodeInstruction(address Word) (Opcode, Word) {
	word := s.store[address]

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

func (s *Ssem) Execute(opcode Opcode, data Word) error {
	switch opcode {
	case JMP:
		s.ci = s.store[data]
	case JRP:
		s.ci += s.store[data]
	case LDN:
		s.a = -s.store[data]
	case STO:
		// this indexing is safe as long as the data extracted earlier (a 5 bit uint for SSEM)
		// is smaller than the number of addresses on the store (32 for SSEM)
		s.store[data] = s.a
	case SUB, SUB2:
		s.a -= s.store[data]
	case CMP:
		if s.a < 0 {
			s.ci += 1
		}
	case STP:
		s.StopFlag = true
	case NUM:
		return fmt.Errorf("encountered unexpected NUM command")
	}

	return nil
}

// Run the machine until STP is encountered or the given amount of cycles is reached.
// Returns the number of cycles executed.
func (s *Ssem) Run(max_cycles uint) (uint, error) {
	var i uint

	// TODO: placeholder, replace this with a proper cycle loop
	for i = 0; i < max_cycles; i++ {
		if err := s.InstructionCycle(); err != nil {
			return i, err
		}
		fmt.Println(s)
		time.Sleep(1 * time.Millisecond)
	}

	return i, nil
}
