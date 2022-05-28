package tracker

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/kydva/jija/peers"
	"github.com/kydva/jija/torrentfile"
	"github.com/zeebo/bencode"
)

type TrackerResponse struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func buildRequestURL(tf *torrentfile.TorrentFile, announce string, peerID [20]byte) (string, error) {
	trackerUrl, err := url.Parse(announce)
	if err != nil {
		return "", err
	}

	query := url.Values{}

	query.Set("info_hash", string(tf.InfoHash[:]))
	query.Set("peer_id", string(peerID[:]))
	query.Set("port", "1337")
	query.Set("uploaded", "0")
	query.Set("downloaded", "0")
	query.Set("compact", "1")
	query.Set("left", strconv.Itoa(tf.Info.Length))

	trackerUrl.RawQuery = query.Encode()

	return trackerUrl.String(), nil
}

func RequestPeers(tf *torrentfile.TorrentFile, peerID [20]byte) ([]peers.Peer, error) {
	announce, err := findTracker(tf)

	if err != nil {
		return nil, err
	}

	url, err := buildRequestURL(tf, announce, peerID)
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

	reader := bencode.NewDecoder(rawResponse.Body)

	err = reader.Decode(&res)

	if err != nil {
		return nil, err
	}

	return peers.Unmarshal([]byte(res.Peers))
}

func findTracker(tf *torrentfile.TorrentFile) (string, error) {
	fmt.Println("Searching available tracker...")

	client := &http.Client{Timeout: 4 * time.Second}

	for _, announces := range tf.AnnounceList {
		for _, announce := range announces {
			if !strings.HasPrefix(announce, "http") {
				continue
			}

			res, err := client.Get(announce)

			if err != nil {
				continue
			}

			defer res.Body.Close()

			fmt.Println("Tracker is found")

			return announce, nil
		}
	}

	return "", fmt.Errorf("there is no available trackers")
}
