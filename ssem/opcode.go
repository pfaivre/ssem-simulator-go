package ssem

import "fmt"

// Represents an operation code, with its binary value
type Opcode uint32

const (
	// Indirect jump
	JMP Opcode = 0b000
	// Relative jump
	JRP Opcode = 0b001
	// Load negative of value from given address to accumulator
	LDN Opcode = 0b010
	// Store accumulator in given address
	STO Opcode = 0b011
	// Substract value in given address from accumulator
	SUB Opcode = 0b100
	// Should not be used. Same effect as SUB
	SUB2 Opcode = 0b101
	// Skip next instruction if accumulator is negative
	CMP Opcode = 0b110
	// Halt the program
	STP Opcode = 0b111
	// Not an instruction. Mnemonic used to set a raw number to the store
	NUM Opcode = 0b1111 // Special differenciation case
)

var mnemonicsMap = map[string]Opcode{
	"JMP":  JMP,
	"JRP":  JRP,
	"LDN":  LDN,
	"STO":  STO,
	"SUB":  SUB,
	"SUB2": SUB2,
	"CMP":  CMP,
	"STP":  STP,
	"NUM":  NUM,
}

var mnemonicsPrint = map[Opcode]string{
	JMP:  "JMP",
	JRP:  "JRP",
	LDN:  "LDN",
	STO:  "STO",
	SUB:  "SUB",
	SUB2: "SUB2",
	CMP:  "CMP",
	STP:  "STP",
	NUM:  "NUM",
}

var mnemonicsNeedOperand = map[Opcode]bool{
	JMP:  true,
	JRP:  true,
	LDN:  true,
	STO:  true,
	SUB:  true,
	SUB2: true,
	CMP:  false,
	STP:  false,
	NUM:  true,
}

func FromString(s string) (Opcode, error) {
	v, presence := mnemonicsMap[s]
	if !presence {
		return 0, fmt.Errorf("unknown mnemonic '%s'", s)
	}
	return v, nil
}

func (o Opcode) String() string {
	return mnemonicsPrint[o]
}

func (o Opcode) NeedsOperand() bool {
	return mnemonicsNeedOperand[o]
}
