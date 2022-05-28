package handshake

import (
	"fmt"
	"io"
)

type Handshake struct {
	Pstr     string
	InfoHash [20]byte
	PeerID   [20]byte
}

const ReservedBytesLen = 8 // Length of reserved bytes that identify supported extensions

func New(infoHash, peerID [20]byte) *Handshake {
	return &Handshake{
		Pstr:     "BitTorrent protocol",
		InfoHash: infoHash,
		PeerID:   peerID,
	}
}

func (h *Handshake) Serialize() []byte {
	buf := make([]byte, len(h.Pstr)+49)
	buf[0] = byte(len(h.Pstr))
	start := 1
	start += copy(buf[start:], h.Pstr)
	start += copy(buf[start:], make([]byte, ReservedBytesLen))
	start += copy(buf[start:], h.InfoHash[:])
	start += copy(buf[start:], h.PeerID[:])
	return buf
}

func Read(r io.Reader) (*Handshake, error) {
	bufLen := make([]byte, 1)
	_, err := io.ReadFull(r, bufLen)
	if err != nil {
		return nil, err
	}

	pstrlen := int(bufLen[0])

	if pstrlen == 0 {
		err := fmt.Errorf("pstrlen cannot be 0")
		return nil, err
	}

	handshakeBuf := make([]byte, 48+pstrlen)
	_, err = io.ReadFull(r, handshakeBuf)
	if err != nil {
		return nil, err
	}

	var infoHash, peerID [20]byte

	copy(infoHash[:], handshakeBuf[pstrlen+ReservedBytesLen:pstrlen+ReservedBytesLen+20])
	copy(peerID[:], handshakeBuf[pstrlen+ReservedBytesLen+20:])

	h := Handshake{
		Pstr:     string(handshakeBuf[0:pstrlen]),
		InfoHash: infoHash,
		PeerID:   peerID,
	}

	return &h, nil
}
