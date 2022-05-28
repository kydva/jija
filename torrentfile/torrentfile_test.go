package torrentfile

import "testing"

func TestOpen(t *testing.T) {
	file, err := Open("../testfile.torrent")
	if err != nil {
		t.Fatal(err)
	}

	expected := TorrentFile{
		Info: TorrentInfo{
			Name:        "Harry_Potter_and_the_Prisoner_of_Azkaban_2004.mkv",
			PieceLength: 1048576,
			Length:      576131414,
		},
	}

	if file.Info.Length != expected.Info.Length {
		t.Fatalf("Info.Length %d != %d", file.Info.Length, expected.Info.Length)
	}

	if file.Info.PieceLength != expected.Info.PieceLength {
		t.Fatalf("Info.PieceLength %d != %d", file.Info.PieceLength, expected.Info.PieceLength)
	}
}

