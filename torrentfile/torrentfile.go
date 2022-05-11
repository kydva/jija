package torrentfile

import (
	"bytes"
	"crypto/sha1"
	"os"

	"github.com/jackpal/bencode-go"
)

const HASH_LENGTH = 20

type hash [HASH_LENGTH]byte

type rawInfo struct {
	Pieces      string
	PieceLength int `bencode:"piece length"`
	Length      int
	Name        string
}

type rawTorrentFile struct {
	Announce string
	Info     rawInfo
}

type TorrentFile struct {
	Announce    string
	InfoHash    hash
	Pieces      []hash
	PieceLength int
	Length      int
	Name        string
}

func Open(path string) (TorrentFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return TorrentFile{}, err
	}
	decoded, err := decode(file)
	if err != nil {
		return TorrentFile{}, err
	}
	return decoded.toTorrentFile()
}

func decode(file *os.File) (rawTorrentFile, error) {
	data := rawTorrentFile{}
	err := bencode.Unmarshal(file, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (info *rawInfo) hash() (hash hash, err error) {
	var buf bytes.Buffer
	err = bencode.Marshal(&buf, *info)
	if err != nil {
		return [20]byte{}, err
	}
	hash = sha1.Sum(buf.Bytes())
	return hash, nil
}

func (info *rawInfo) splitPieces() ([]hash, error) {
	buf := []byte(info.Pieces)
	piecesCount := len(buf) / HASH_LENGTH
	pieces := make([]hash, piecesCount)

	for i := 0; i < piecesCount; i++ {
		start := i * HASH_LENGTH
		end := (i + 1) * HASH_LENGTH
		copy(pieces[i][:], buf[start:end])
	}
	return pieces, nil
}

func (file *rawTorrentFile) toTorrentFile() (TorrentFile, error) {
	infoHash, err := file.Info.hash()
	if err != nil {
		return TorrentFile{}, err
	}

	pieces, err := file.Info.splitPieces()
	if err != nil {
		return TorrentFile{}, err
	}

	torrentFile := TorrentFile{
		Announce:    file.Announce,
		InfoHash:    infoHash,
		Pieces:      pieces,
		PieceLength: file.Info.PieceLength,
		Length:      file.Info.Length,
		Name:        file.Info.Name,
	}

	return torrentFile, nil
}
