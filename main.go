package main

import (
	"fmt"

	"github.com/concurrency-8/torrent"
)

func main() {

	torrent.DownloadFromFile("/home/shraiysh/Downloads/CBA8DE1D68607F71B1DBBDAB6319C9B3257A8E83.torrent", 6881)
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
