package ssem

import (
	"fmt"
	"math/bits"
	"strings"
)

// A word is represented by a 32 bits integer
type Word int32

// Contains the machine's memory
type Store [WORD_COUNT]Word

func (s Store) String() string {
	builder := strings.Builder{}
	for _, w := range s {
		builder.WriteString(fmt.Sprintf("%032b\n", bits.Reverse32(uint32(w))))
	}

	replacer := strings.NewReplacer("0", ".", "1", "#")
	return replacer.Replace(builder.String())
}
