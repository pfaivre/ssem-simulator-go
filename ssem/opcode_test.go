package ssem

import "testing"

func TestFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Opcode
		wantErr bool
	}{
		{"JMP", "JMP", JMP, false},
		{"JRP", "JRP", JRP, false},
		{"LDN", "LDN", LDN, false},
		{"STO", "STO", STO, false},
		{"SUB", "SUB", SUB, false},
		{"SUB2", "SUB2", SUB2, false},
		{"CMP", "CMP", CMP, false},
		{"STP", "STP", STP, false},
		{"NUM", "NUM", NUM, false},
		{"unknown mnemonic", "XYZ", 0, true},
		{"lower case mnemonic", "jmp", 0, true},
		{"empty mnemonic", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromString(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("FromString(%q) should return an error", tt.input)
				}
				return
			}
			if err != nil {
				t.Errorf("FromString(%q) returned an unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Errorf("FromString(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestOpcodeString(t *testing.T) {
	tests := []struct {
		opcode Opcode
		want   string
	}{
		{JMP, "JMP"},
		{JRP, "JRP"},
		{LDN, "LDN"},
		{STO, "STO"},
		{SUB, "SUB"},
		{SUB2, "SUB2"},
		{CMP, "CMP"},
		{STP, "STP"},
		{NUM, "NUM"},
	}

	for _, tt := range tests {
		if got := tt.opcode.String(); got != tt.want {
			t.Errorf("Opcode(%d).String() = %q, want %q", tt.opcode, got, tt.want)
		}
	}
}

func TestNeedsOperand(t *testing.T) {
	tests := []struct {
		opcode Opcode
		want   bool
	}{
		{JMP, true},
		{JRP, true},
		{LDN, true},
		{STO, true},
		{SUB, true},
		{SUB2, true},
		{CMP, false},
		{STP, false},
		{NUM, true},
	}

	for _, tt := range tests {
		if got := tt.opcode.NeedsOperand(); got != tt.want {
			t.Errorf("Opcode(%d).NeedsOperand() = %t, want %t", tt.opcode, got, tt.want)
		}
	}
}
