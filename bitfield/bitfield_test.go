package bitfield

import (
	"bytes"
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

func TestSetPiece(t *testing.T) {
	bf := &Bitfield{0b11011011, 0b11100101}

	indexes := []int{2, 3, 11}

	expectedOutput := &Bitfield{0b11111011, 0b11110101}

	for _, index := range indexes {
		bf.SetPiece(index)
	}

	if !bytes.Equal(*bf, *expectedOutput) {
		t.Fatalf("% 08b != % 08b", bf, expectedOutput)
	}
}
