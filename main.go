package main

import (
	"fmt"

	"github.com/concurrency-8/torrent"
)

func main() {

	torrent.DownloadFromFile("/home/shraiysh/go/src/github.com/concurrency-8/test_torrents/52 MB.torrent", 6881)
	// torrentfile, _ := parser.ParseFromFile("./test_torrents/ubuntu.iso.torrent")
	// fmt.Println(torrentfile.PieceLength)
	// u, err := url.Parse(torrentfile.Announce[0])
	// fmt.Println(u)
	// if err != nil {
	// 	return
	// }

	// clientReport := tracker.GetClientStatusReport(torrentfile, 6881)

	// announce, err := tracker.GetPeers(u, clientReport)

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// pieces := piece.NewPieceTracker(torrentfile)
	// torrent.Download(announce.Peers[1], clientReport, pieces)
	fmt.Println("It's over")
}
