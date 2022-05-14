package tracker

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/jackpal/bencode-go"
	"github.com/kydva/jija/peers"
	"github.com/kydva/jija/torrentfile"
)

type TrackerResponse struct {
	Interval string
	Peers    string
}

func buildTrackerURL(tf *torrentfile.TorrentFile, peerID string) (string, error) {
	trackerUrl, err := url.Parse(tf.Announce)
	if err != nil {
		return "", err
	}

	infohash, err := tf.Info.Hash()
	if err != nil {
		return "", err
	}

	query := url.Values{}

	query.Set("info_hash", infohash)
	query.Set("peer_id", peerID)
	query.Set("port", "1337")
	query.Set("uploaded", "0")
	query.Set("downloaded", "0")
	query.Set("compact", "1")
	query.Set("left", strconv.Itoa(tf.Info.Length))

	trackerUrl.RawQuery = query.Encode()

	return trackerUrl.String(), nil
}

func RequestPeers(tf *torrentfile.TorrentFile, peerID string) ([]peers.Peer, error) {
	url, err := buildTrackerURL(tf, peerID)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 15 * time.Second}

	rawResponse, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer rawResponse.Body.Close()

	res := TrackerResponse{}
	err = bencode.Unmarshal(rawResponse.Body, &res)
	if err != nil {
		return nil, err
	}

	return peers.Unmarshal([]byte(res.Peers))
}
