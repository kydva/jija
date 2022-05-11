package torrentfile

import (
	"testing"
)

func TestOpen(t *testing.T) {
	data, err := Open("testdata/file.torrent")
	if err != nil {
		t.Fatal(err)
	}

	expected := TorrentFile{
		Announce:    "udp://tracker.opentrackr.org:1337/announce",
		PieceLength: 262144,
		Length:      363548672,
	}

	if data.Announce != expected.Announce {
		t.Fatalf("Announce %s != %s", data.Announce, expected.Announce)
	}

	if data.Length != expected.Length {
		t.Fatalf("Length %d != %d", data.Length, expected.Length)
	}

	if data.PieceLength != expected.PieceLength {
		t.Fatalf("PieceLength %d != %d", data.PieceLength, expected.PieceLength)
	}

	fifthPieceHash := [20]byte{46, 0, 15, 167, 232, 87, 89, 199, 244, 194, 84, 212, 217, 195, 62, 244, 129, 228, 89, 167}

	if data.Pieces[5] != fifthPieceHash {
		t.Fatalf("5th piece hash %x != %x", data.Pieces[5], fifthPieceHash)
	}
}
