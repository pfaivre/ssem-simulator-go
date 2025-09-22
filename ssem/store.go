package ssem

import (
	"fmt"
	"strings"
)

// A word is represented by a 32 bits integer
type Word int32

// Contains the machine's memory
type Store [WORD_COUNT]Word

func (s Store) String() string {
	builder := strings.Builder{}
	builder.Grow(1500) // takes approximately 1056 bytes to print the store
	for _, w := range s {
		AppendBinary(&builder, w)
		builder.WriteString("\n")
	}

	return builder.String()
}

// var replacer = strings.NewReplacer("0", ".", "1", "#")

func AppendBinary(b *strings.Builder, w Word) {
	// b.WriteString(
	// replacer.Replace(
	// Reverse(fmt.Sprintf("%032b", uint32(w)))))
	b.WriteString(Reverse(fmt.Sprintf("%032b", uint32(w))))
}

// Returns a new string with each rune in reverse order
func Reverse(s string) string {
	// TODO: consider swapping slices of bytes instead for more efficiency, if possible
	// func Reverse(s []byte)
	var b strings.Builder
	b.Grow(len(s))
	runes := []rune(s)
	for i := len(runes) - 1; i >= 0; i-- {
		b.WriteRune(runes[i])
	}
	return b.String()
}
