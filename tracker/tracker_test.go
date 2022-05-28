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

	expectedUrl := "http://tracker.opentrackr.org/announce?compact=1&downloaded=0&info_hash=%C0%F9o%E1%D7%F1%CCl3%DC23a%5B%1E%FBv%DCKp&left=576131414&peer_id=%00%01%02%03%04%05%06%07%08%09%0A%0B%0C%0D%0E%0F%10%11%12%13&port=1337&uploaded=0"
	url, err := buildRequestURL(&tf, "http://tracker.opentrackr.org/announce", [20]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19})
	if err != nil {
		t.Fatal(err)
	}

	if url != expectedUrl {
		t.Fatalf("%s != %s", url, expectedUrl)
	}
}

// TODO: Test RequestPeers() func
// func TestRequestPeers(t *testing.T) { }
