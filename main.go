package main

import (
	"fmt"
	"github.com/concurrency-8/parser"
)

func main() {
	torrentfile, _ := parser.ParseFromFile("./test_torrents/[TorrentCounter.to].Thugs.of.Hindostan.2018.Hindi.720p.BluRay.x264.[1.2GB].[MP4].torrent")
	for _, A := range torrentfile.Files {
		fmt.Println(A.Path)
	}
	fmt.Println(torrentfile.Length)
	fmt.Println(torrentfile.InfoHash)

}
