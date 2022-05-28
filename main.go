package main

import (
	"log"
	"os"

	"github.com/kydva/jija/download"
	"github.com/kydva/jija/torrentfile"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("You must specify torrent file")
	}

	torrentPath := os.Args[1]

	var downloadPath string

	if len(os.Args) > 2 {
		downloadPath = os.Args[2]
	}

	tf, err := torrentfile.Open(torrentPath)
	if err != nil {
		log.Fatal(err)
	}

	err = download.Download(&tf, downloadPath)

	if err != nil {
		log.Fatal(err)
	}
}
