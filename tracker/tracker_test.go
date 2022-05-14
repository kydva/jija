package tracker

import (
	"testing"

	"github.com/kydva/jija/torrentfile"
)

func TestBuildTrackerUrl(t *testing.T) {
	tf, err := torrentfile.Open("../torrentfile/testdata/file.torrent")
	if err != nil {
		t.Fatal(err)
	}

	expectedUrl := "udp://tracker.opentrackr.org:1337/announce?compact=1&downloaded=0&info_hash=%8BDW%BC%F0%ABV%A8%B5%08%DE%BA%EC%B5%8A%A7%5D%E3%CC%2C&left=363548672&peer_id=xxxxxxxxxxxxxxx&port=1337&uploaded=0"
	url, err := buildTrackerURL(&tf, "xxxxxxxxxxxxxxx")
	if err != nil {
		t.Fatal(err)
	}

	if url != expectedUrl {
		t.Fatalf("%s != %s", url, expectedUrl)
	}
}
