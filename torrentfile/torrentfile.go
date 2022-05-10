package torrentfile

import (
	"os"

	"github.com/jackpal/bencode-go"
)

type TorrentInfo struct {
	Pieces      string
	PieceLength int `bencode:"piece length"`
	Length      int
	Name        string
}

type TorrentFile struct {
	Announce string
	Info     TorrentInfo
}

func Open(path string) (TorrentFile, error) {
	file, err := os.Open(path)
	decoded := TorrentFile{}
	if err != nil {
		return decoded, err
	}
	err = bencode.Unmarshal(file, &decoded)
	if err != nil {
		return decoded, err
	}
	return decoded, nil
}
