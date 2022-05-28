package download

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/kydva/jija/connection"
	"github.com/kydva/jija/message"
	"github.com/kydva/jija/peers"
	"github.com/kydva/jija/torrentfile"
	"github.com/kydva/jija/tracker"
)

const MaxBlockSize = 16384 // the largest number of bytes a request can ask for

const MaxBacklog = 5 // the largest number of unfulfilled requests a client can make

type piece struct {
	index  int
	hash   [20]byte
	length int
}

type pieceResult struct {
	index int
	buf   []byte
}

type pieceProgress struct {
	index      int
	conn       *connection.Connection
	buf        []byte
	downloaded int
	requested  int
	backlog    int
}

func handleMessage(state *pieceProgress) error {
	msg, err := state.conn.Read()
	if err != nil {
		return err
	}

	// keep-alive
	if msg == nil {
		return nil
	}

	switch msg.ID {
	case message.MsgUnchoke:
		state.conn.Choked = false
	case message.MsgChoke:
		state.conn.Choked = true
	case message.MsgHave:
		index, err := message.ParseHave(msg)
		if err != nil {
			return err
		}
		state.conn.Bitfield.SetPiece(index)
	case message.MsgPiece:
		n, err := message.ParsePiece(state.index, state.buf, msg)
		if err != nil {
			return err
		}
		state.downloaded += n
		state.backlog--
	}
	return nil
}

func downloadPiece(c *connection.Connection, piece *piece) ([]byte, error) {
	state := pieceProgress{
		index: piece.index,
		conn:  c,
		buf:   make([]byte, piece.length),
	}

	c.SetDeadline(time.Now().Add(30 * time.Second))
	defer c.SetDeadline(time.Time{})

	for state.downloaded < piece.length {
		if !state.conn.Choked {
			for state.backlog < MaxBacklog && state.requested < piece.length {
				blockSize := MaxBlockSize

				// Last block can be shorter than others
				if piece.length-state.requested < blockSize {
					blockSize = piece.length - state.requested
				}

				err := c.SendRequest(piece.index, state.requested, blockSize)
				if err != nil {
					return nil, err
				}

				state.backlog++
				state.requested += blockSize
			}
		}

		err := handleMessage(&state)
		if err != nil {
			return nil, err
		}
	}

	return state.buf, nil
}

func (piece *piece) checkIntegrity(buf []byte) error {
	hash := sha1.Sum(buf)
	if !bytes.Equal(hash[:], piece.hash[:]) {
		return fmt.Errorf("index %d failed integrity check", piece.index)
	}
	return nil
}

func startDownloadWorker(peerID [20]byte, infoHash [20]byte, peer peers.Peer, workQueue chan *piece, results chan *pieceResult) {
	conn, err := connection.New(peer, peerID, infoHash)

	if err != nil {
		fmt.Printf("❌ - %s\n", peer.IP)
		return
	}

	fmt.Printf("✅ - %s\n", peer.IP)

	defer conn.Close()
	conn.SendUnchoke()
	conn.SendInterested()

	for piece := range workQueue {
		if !conn.Bitfield.HasPiece(piece.index) {
			workQueue <- piece
			continue
		}

		buf, err := downloadPiece(conn, piece)
		if err != nil {
			workQueue <- piece
			return
		}

		err = piece.checkIntegrity(buf)
		if err != nil {
			fmt.Printf("Piece #%d failed integrity check\n", piece.index)
			workQueue <- piece
			continue
		}

		conn.SendHave(piece.index)
		results <- &pieceResult{piece.index, buf}
	}
}

func Download(tf *torrentfile.TorrentFile, dir string) error {
	var peerID [20]byte
	rand.Read(peerID[:])

	peers, err := tracker.RequestPeers(tf, peerID)
	if err != nil {
		return err
	}

	workQueue := make(chan *piece, len(tf.Info.Pieces))
	results := make(chan *pieceResult)

	pieces, err := tf.SplitPieces()

	if err != nil {
		return err
	}

	for index, hash := range pieces {
		length := tf.CalculatePieceSize(index)
		workQueue <- &piece{index, hash, length}
	}

	fmt.Println("Peers list:")
	for _, peer := range peers {
		go startDownloadWorker(peerID, tf.InfoHash, peer, workQueue, results)
	}

	buf := make([]byte, tf.Info.Length)

	donePieces := 0

	for donePieces < len(pieces) {
		res := <-results
		begin, end := tf.CalculateBoundsForPiece(res.index)
		copy(buf[begin:end], res.buf)
		donePieces++

		percent := float64(donePieces) / float64(len(pieces)) * 100
		fmt.Printf("%0.2f%% Done \n", percent)
	}

	close(workQueue)

	path := filepath.Join(dir, tf.Info.Name)

	outFile, err := os.Create(path)
	if err != nil {
		return err
	}

	defer outFile.Close()
	_, err = outFile.Write(buf)
	if err != nil {
		return err
	}

	return nil
}
