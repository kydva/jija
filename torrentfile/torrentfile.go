package torrentfile

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"os"

	"github.com/zeebo/bencode"
)

type TorrentInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type TorrentFile struct {
	AnnounceList [][]string  `bencode:"announce-list"`
	Info         TorrentInfo `bencode:"info"`
	InfoHash     [20]byte
}

func Open(path string) (TorrentFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return TorrentFile{}, err
	}

	tf := TorrentFile{}

	decoder := bencode.NewDecoder(file)

	decoder.Decode(&tf)

	tf.InfoHash, err = tf.Info.hash()
	if err != nil {
		return tf, err
	}

	fmt.Println(tf.AnnounceList)

	return tf, nil
}

func (info *TorrentInfo) hash() ([20]byte, error) {
	var buf bytes.Buffer

	encoder := bencode.NewEncoder(&buf)

	encoder.Encode(*info)

	hash := sha1.Sum(buf.Bytes())
	return hash, nil
}

func (tf *TorrentFile) SplitPieces() ([][20]byte, error) {
	const hashLen = 20
	buf := []byte(tf.Info.Pieces)
	piecesCount := len(buf) / hashLen
	pieces := make([][20]byte, piecesCount)

	for i := 0; i < piecesCount; i++ {
		start := i * hashLen
		end := (i + 1) * hashLen
		copy(pieces[i][:], buf[start:end])
	}
	return pieces, nil
}

func (tf *TorrentFile) CalculateBoundsForPiece(index int) (begin int, end int) {
	begin = index * tf.Info.PieceLength
	end = begin + tf.Info.PieceLength
	if end > tf.Info.Length {
		end = tf.Info.Length
	}
	return begin, end
}

func (tf *TorrentFile) CalculatePieceSize(index int) int {
	begin, end := tf.CalculateBoundsForPiece(index)
	return end - begin
}
