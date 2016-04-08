package main

import (
	"fmt"
	"github.com/mrd0ll4r/poke"
	"github.com/mrd0ll4r/poke/http"
	"log"
	"net"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	req := poke.AnnounceRequest{
		InfoHash: poke.InfoHash("8ff5ee1ddd3ddf101d1733d81512dc283ea8ef2f"),
		Peer: poke.Peer{
			IP:   net.ParseIP("12.3.4.5"),
			Port: 1235,
			ID:   "-ADMN64-000000000000"},
		Event:   poke.EventStarted,
		Numwant: 50,
		Compact: false,
	}

	c, err := http.NewClient("http://localhost:6882/announce")
	if err != nil {
		log.Fatal(err)
	}

	resp, err := c.Announce(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", resp)
}
