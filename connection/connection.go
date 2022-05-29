package connection

import (
	"bytes"
	"fmt"
	"net"
	"time"

	"github.com/kydva/jija/bitfield"
	"github.com/kydva/jija/handshake"
	"github.com/kydva/jija/message"
	"github.com/kydva/jija/peers"
)

// P2P connection
type Connection struct {
	net.Conn
	Choked   bool
	Bitfield bitfield.Bitfield
	peer     peers.Peer
	infoHash [20]byte
	peerID   [20]byte
}

func completeHandshake(conn net.Conn, infohash, peerID [20]byte) error {
	conn.SetDeadline(time.Now().Add(3 * time.Second))
	defer conn.SetDeadline(time.Time{})

	req := handshake.New(infohash, peerID)

	_, err := conn.Write(req.Serialize())
	if err != nil {
		return nil
	}

	res, err := handshake.Read(conn)
	if err != nil {
		return nil
	}

	if !bytes.Equal(res.InfoHash[:], infohash[:]) {
		return fmt.Errorf("expected infohash %x but got %x", res.InfoHash, infohash)
	}

	return nil
}

func recieveBitfield(conn net.Conn) (bitfield.Bitfield, error) {
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetDeadline(time.Time{})

	msg, err := message.Read(conn)
	if err != nil {
		return nil, err
	}
	if msg == nil {
		err := fmt.Errorf("expected bitfield but got nothing")
		return nil, err
	}
	if msg.ID != message.MsgBitfield {
		err := fmt.Errorf("expected bitfield but got ID %d", msg.ID)
		return nil, err
	}

	return msg.Payload, nil
}

func New(peer peers.Peer, peerID, infoHash [20]byte) (*Connection, error) {
	conn, err := net.DialTimeout("tcp", peer.String(), 3*time.Second)
	if err != nil {
		return nil, err
	}

	err = completeHandshake(conn, infoHash, peerID)
	if err != nil {
		conn.Close()
		return nil, err
	}

	bf, err := recieveBitfield(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Connection{
		Conn:     conn,
		Choked:   true,
		Bitfield: bf,
		peer:     peer,
		infoHash: infoHash,
		peerID:   peerID,
	}, nil
}

func (c *Connection) Read() (*message.Message, error) {
	msg, err := message.Read(c.Conn)
	return msg, err
}

func (c *Connection) SendRequest(index, begin, length int) error {
	req := message.NewRequest(index, begin, length)
	_, err := c.Write(req.Serialize())
	return err
}

func (c *Connection) SendInterested() error {
	msg := message.Message{ID: message.MsgInterested}
	_, err := c.Write(msg.Serialize())
	return err
}

func (c *Connection) SendUnchoke() error {
	msg := message.Message{ID: message.MsgUnchoke}
	_, err := c.Write(msg.Serialize())
	return err
}

func (c *Connection) SendHave(index int) error {
	msg := message.NewHave(index)
	_, err := c.Write(msg.Serialize())
	return err
}
