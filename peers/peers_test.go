package peers

import (
	"net"
	"reflect"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	input := []byte{127, 0, 0, 1, 0x1A, 0xE1, 1, 1, 1, 1, 0x04, 0x1F}
	expectedOutput := []Peer{
		{IP: net.IP{127, 0, 0, 1}, Port: 6881}, // [0x1A, 0xE1] = 0x1AE1 = 6881
		{IP: net.IP{1, 1, 1, 1}, Port: 1055},   // [0x04, 0x1F] = 0x041F = 1055
	}

	output, err := Unmarshal(input)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(output, expectedOutput) {
		t.Fatalf("%+v\n != %+v\n", output, expectedOutput)
	}
}
