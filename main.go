package main

import (
	"fmt"
	"bufio"
	"os"
	"io"
	"strings"
	// "github.com/concurrency-8/parser"
	// "github.com/concurrency-8/piece"
	"github.com/concurrency-8/torrent"
	// "github.com/concurrency-8/tracker"
	// "net/url"
)

func main() {
	fmt.Printf("GoTorrent v1.0\nType \"help\" or \"license\" for more information.\n")
	scanner := bufio.NewReader(os.Stdin)
	ports := make([]int, 1)
	ports[0] = 7000
	var text string
	var err error
	text = ""
	for text!="exit" && err!=io.EOF{
		fmt.Printf(">>> ")
		text, err = scanner.ReadString('\n')
		switch {
		case text=="help\n":
			fmt.Println("\n\tUsage:")
			fmt.Println("\t\tfile [path] -- Downloads the torrent file specified by path.")
			fmt.Println("\t\tshow   -- Shows current files that are being downloaded.")
			fmt.Println("\t\tpause [id]  -- Pauses the download of torrent specified by torrent index aka id.\n\t\t\t To show torrent index, use 'show'.")
			fmt.Println("\t\tresume [id] -- Resumes the download of torrent specified by torrent index aka id.\n\t\t\t To show torrent index, use 'show'.")
			fmt.Println("\t\texit -- Exit shell.")
		case len(text)>4 && text[:4]=="file":
			text := strings.TrimSpace(text[4:])
			fmt.Println(text)
			//Get the next port to download, maynot be available, download will fail with error.
			next := ports[len(ports)-1] + 1
			ports = append(ports, next)
			go torrent.DownloadFromFile(text, next)
		
		}

	}
}
