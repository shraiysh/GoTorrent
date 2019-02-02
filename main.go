package main

import (
	"fmt"
	"github.com/concurrency-8/parser"
	"github.com/concurrency-8/tracker"
	"net/url"
)

func main() {
	torrentfile, _ := parser.ParseFromFile("./test_torrents/ubuntu.iso.torrent")
	u, err := url.Parse(torrentfile.Announce[0])
	fmt.Println(u)
	if err != nil {
		return
	}

	clientReport := tracker.GetClientStatusReport(torrentfile, 6881)

	announce, err := tracker.GetPeers(u, clientReport)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(announce.Peers)

}
