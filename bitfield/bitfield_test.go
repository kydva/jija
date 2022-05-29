package bitfield

import (
	"testing"
)

func TestHasPiece(t *testing.T) {
	bf := &Bitfield{0b11011011, 0b11100101}

	tests := map[int]bool{
		2:  false,
		3:  true,
		11: false,
	}

	for index, output := range tests {
		if bf.HasPiece(index) != output {
			t.Fatalf("Failed check for piece with index %d", index)
		}
	}

}
