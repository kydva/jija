package torrentfile

import "testing"

func TestOpen(t *testing.T) {
	file, err := Open("testdata/file.torrent")
	if err != nil {
		t.Fatal(err)
	}

	expected := TorrentFile{
		Announce: "udp://tracker.openbittorrent.com:80",
		Info: TorrentInfo{
			PieceLength: 65536,
			Length:      20,
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
