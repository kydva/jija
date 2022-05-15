package tracker

import (
	"testing"

	"github.com/kydva/jija/torrentfile"
)

func TestBuildTrackerUrl(t *testing.T) {
	tf, err := torrentfile.Open("../testfile.torrent")
	if err != nil {
		t.Fatal(err)
	}

	expectedUrl := "http://tracker.archlinux.org:6969/announce?compact=1&downloaded=0&info_hash=%87%0B%B5%D8%08%B2%C9R%11PA%EA%82k%1A%5Dc%1F%EBD&left=670040064&peer_id=xxxxxxxxxxxxxxx&port=1337&uploaded=0"
	url, err := buildTrackerURL(&tf, "xxxxxxxxxxxxxxx")
	if err != nil {
		t.Fatal(err)
	}

	if url != expectedUrl {
		t.Fatalf("%s != %s", url, expectedUrl)
	}
}

// TODO: Test RequestPeers() func
// func TestRequestPeers(t *testing.T) { }
