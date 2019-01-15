package main

import ( 
			"fmt"
			"./parser"
			"./tracker"
		)

func main() {
	torrent_file,_ := parser.ParseFromFile("./test_torrents/[TorrentCounter.to].Thugs.of.Hindostan.2018.Hindi.720p.BluRay.x264.[1.2GB].[MP4].torrent")
	fmt.Println(torrent_file)

}