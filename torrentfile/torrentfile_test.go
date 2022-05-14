package torrentfile

import "testing"

func TestOpen(t *testing.T) {
	file, err := Open("testdata/file.torrent")
	if err != nil {
		t.Fatal(err)
	}

	expected := TorrentFile{
		Announce: "udp://tracker.opentrackr.org:1337/announce",
		Info: TorrentInfo{
			PieceLength: 262144,
			Length:      363548672,
		},
	}

	if file.Announce != expected.Announce {
		t.Fatalf("Announce %s != %s", file.Announce, expected.Announce)
	}

	if file.Info.Length != expected.Info.Length {
		t.Fatalf("Info.Length %d != %d", file.Info.Length, expected.Info.Length)
	}

	if file.Info.PieceLength != expected.Info.PieceLength {
		t.Fatalf("Info.PieceLength %d != %d", file.Info.PieceLength, expected.Info.PieceLength)
	}
}

// fifthPieceHash := [20]byte{46, 0, 15, 167, 232, 87, 89, 199, 244, 194, 84, 212, 217, 195, 62, 244, 129, 228, 89, 167}

// if data.Pieces[5] != fifthPieceHash {
// 	t.Fatalf("5th piece hash %x != %x", data.Pieces[5], fifthPieceHash)
// }
