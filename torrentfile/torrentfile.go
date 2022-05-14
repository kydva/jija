package torrentfile

import (
	"bytes"
	"crypto/sha1"
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
	if err != nil {
		return TorrentFile{}, err
	}

	decoded := TorrentFile{}
	err = bencode.Unmarshal(file, &decoded)
	if err != nil {
		return decoded, err
	}

	return decoded, nil
}

func (info *TorrentInfo) Hash() (string, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *info)
	if err != nil {
		return "", err
	}
	hash := sha1.Sum(buf.Bytes())
	return string(hash[:]), nil
}

func (info *TorrentInfo) SplitPieces() ([][20]byte, error) {
	const hashLen = 20
	buf := []byte(info.Pieces)
	piecesCount := len(buf) / hashLen
	pieces := make([][20]byte, piecesCount)

	for i := 0; i < piecesCount; i++ {
		start := i * hashLen
		end := (i + 1) * hashLen
		copy(pieces[i][:], buf[start:end])
	}
	return pieces, nil
}
