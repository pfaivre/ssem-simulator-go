package ssem

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Writes the given content to a temporary file and returns its path
func writeTempFile(t *testing.T, name string, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestReadAsm(t *testing.T) {
	content := `; Test program

00 NUM 1    ; Raw value
01 NUM -1   ; Negative raw value
02 LDN 31   ; Load negative of counter
03 SUB 0
04 STO 1
05 CMP
06 STP
`

	s := NewSsem()
	s.ci = 7
	s.a = 42

	err := s.ReadAsm(writeTempFile(t, "program.asm", content))
	if err != nil {
		t.Fatalf("ReadAsm returned an unexpected error: %v", err)
	}

	expected := Store{}
	expected[0] = 1
	expected[1] = -1
	expected[2] = Word(LDN)<<OPCODE_START | 31
	expected[3] = Word(SUB) << OPCODE_START
	expected[4] = Word(STO)<<OPCODE_START | 1
	expected[5] = Word(CMP) << OPCODE_START
	expected[6] = Word(STP) << OPCODE_START

	if s.store != expected {
		t.Errorf("store does not match the assembled program")
	}
	if s.ci != 0 {
		t.Errorf("ci = %d, want 0", s.ci)
	}
	if s.a != 0 {
		t.Errorf("a = %d, want 0", s.a)
	}
	if s.IsStopped() {
		t.Errorf("machine should be started after loading")
	}
}

func TestReadAsmErrors(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr string
	}{
		{
			"bad index",
			"xx LDN 1\n",
			"line 1: unable to parse numeric value 'xx'",
		},
		{
			"negative index",
			"-1 LDN 1\n",
			"line 1 index out of bound",
		},
		{
			"index out of bound",
			"32 JMP 0\n",
			"line 1 index out of bound",
		},
		{
			"unknown mnemonic",
			"00 XYZ 1\n",
			"line 1: unknown mnemonic 'XYZ'",
		},
		{
			"missing operand",
			"00 JMP\n",
			"line 1: mnemonic 'JMP' needs an operand",
		},
		{
			"NUM without operand",
			"00 NUM\n",
			"line 1: mnemonic 'NUM' needs an operand",
		},
		{
			"operand on STP",
			"00 STP 1\n",
			"line 1: mnemonic 'STP' does not accept an operand",
		},
		{
			"bad operand",
			"00 JMP xx\n",
			"line 1: unable to parse numeric value 'xx'",
		},
		{
			"error on the correct line",
			"; comment\n\n00 NUM 1\n01 XYZ 1\n",
			"line 4: unknown mnemonic 'XYZ'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSsem()
			err := s.ReadAsm(writeTempFile(t, "program.asm", tt.content))
			if err == nil {
				t.Fatalf("ReadAsm should return an error")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("ReadAsm error = %q, want it to contain %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestReadAsmMissingFile(t *testing.T) {
	s := NewSsem()
	err := s.ReadAsm(filepath.Join(t.TempDir(), "does-not-exist.asm"))
	if err == nil {
		t.Errorf("ReadAsm on a missing file should return an error")
	}
}

func TestReadSnp(t *testing.T) {
	content := `; SSEM snapshot

00: 10000000000000000000000000000000
01: 01000000000000000000000000000000
02: 11111111111111111111111111111111
31: 00000000000000000000000000000000
`

	s := NewSsem()
	s.ci = 7
	s.a = 42

	err := s.ReadSnp(writeTempFile(t, "snapshot.snp", content))
	if err != nil {
		t.Fatalf("ReadSnp returned an unexpected error: %v", err)
	}

	expected := Store{}
	expected[0] = 1
	expected[1] = 2
	expected[2] = -1

	if s.store != expected {
		t.Errorf("store does not match the snapshot")
	}
	if s.ci != 0 {
		t.Errorf("ci = %d, want 0", s.ci)
	}
	if s.a != 0 {
		t.Errorf("a = %d, want 0", s.a)
	}
	if s.IsStopped() {
		t.Errorf("machine should be started after loading")
	}
}

func TestReadSnpErrors(t *testing.T) {
	WORD := "00000000000000000000000000000000"

	tests := []struct {
		name    string
		content string
		wantErr string
	}{
		{
			"bad index",
			"xx: " + WORD + "\n",
			"line 1: unable to parse numeric value 'xx'",
		},
		{
			"index out of bound",
			"32: " + WORD + "\n",
			"line 1: index out of bound",
		},
		{
			"missing word without separator",
			"00\n",
			"line 1: missing word for index 0",
		},
		{
			"missing word after separator",
			"00:\n",
			"line 1: missing word for index 0",
		},
		{
			"bad binary value",
			"00: 2\n",
			"line 1: unable to parse binary value '2'",
		},
		{
			"error on the correct line",
			"00: " + WORD + "\n01\n",
			"line 2: missing word for index 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSsem()
			err := s.ReadSnp(writeTempFile(t, "snapshot.snp", tt.content))
			if err == nil {
				t.Fatalf("ReadSnp should return an error")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("ReadSnp error = %q, want it to contain %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestReadSnpMissingFile(t *testing.T) {
	s := NewSsem()
	err := s.ReadSnp(filepath.Join(t.TempDir(), "does-not-exist.snp"))
	if err == nil {
		t.Errorf("ReadSnp on a missing file should return an error")
	}
}

// Loads the sample files shipped with the repository
func TestLoadSamples(t *testing.T) {
	s := NewSsem()
	if err := s.ReadAsm(filepath.Join("..", "samples", "ssem", "fibonacci.asm")); err != nil {
		t.Errorf("ReadAsm on fibonacci.asm returned an unexpected error: %v", err)
	}
	if s.store == (Store{}) {
		t.Errorf("fibonacci.asm loaded an empty store")
	}

	if err := s.ReadSnp(filepath.Join("..", "samples", "ssem", "nightmare.snp")); err != nil {
		t.Errorf("ReadSnp on nightmare.snp returned an unexpected error: %v", err)
	}
	if s.store == (Store{}) {
		t.Errorf("nightmare.snp loaded an empty store")
	}
}
