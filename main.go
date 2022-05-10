package main

import (
	"log"
	"os"

	"github.com/kydva/jija/torrentfile"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Torrent file is not specified")
	}
	path := os.Args[1]
	_, err := torrentfile.Open(path)
	if err != nil {
		log.Fatal(err)
	}
}
