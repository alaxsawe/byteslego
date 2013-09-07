package main

import (
	"flag"
	"fmt"

	"github.com/marconi/byteslego/torrent"
	"github.com/marconi/byteslego/tracker"
)

const PORT = 5678

var filename string

func init() {
	flag.StringVar(&filename, "filename", "", "Path to .torrent file.")
}

func main() {
	flag.Parse()
	t := torrent.New(filename)
	tkrs := t.GetTrackers()

	for _, tkr := range tkrs {
		r, err := tracker.NewRequest(tkr.Announce,
			t.InfoHash,
			t.PeerId,
			PORT,
			tracker.STARTED)
		if err != nil {
			fmt.Println(err, "\n")
			continue
		}
		fmt.Println("Trying: ", tkr.Announce)
		if _, err := tkr.GetPeers(r); err != nil {
			fmt.Println(err, "\n")
		}
	}
}
