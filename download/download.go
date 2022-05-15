package download

import (
	"crypto/rand"

	"github.com/kydva/jija/torrentfile"
	"github.com/kydva/jija/tracker"
)

func Download(t *torrentfile.TorrentFile) error {
	var peerID [20]byte
	rand.Read(peerID[:])

	peers, err := tracker.RequestPeers(t, string(peerID[:]))
	if err != nil {
		return err
	}

	
}
