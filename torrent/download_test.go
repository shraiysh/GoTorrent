package torrent

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/concurrency-8/parser"
	"github.com/concurrency-8/tracker"
	"net/url"
	"math/rand"
	"net"
	"testing"
)

// TestOnWholeMessage tests torrent/download.go : onWholeMessage(*kwargs)
func TestOnWholeMessage(t *testing.T) {
	fmt.Println("Testing torrent/download.go : onWholeMessage(*kwargs)")

	//length := rand.Intn(100) // generate random handshake message
	length := 220
	message := make([]byte, length+49)
	rand.Read(message)
	message[0] = uint8(length)

	client, server := net.Pipe() // create a client and server connection

	go func() {
		server.Write(message)
		server.Close() // close after writing out all data
	}()

	err := onWholeMessage(client, func(b []byte, client net.Conn) error{ // mock Message Handler
		assert.Equal(t, len(b), int(b[0])+49, "length not equal")
		assert.Equal(t, b, message, "message received not same")
		return assert.AnError
	})
	assert.Equal(t, err, fmt.Errorf("EOF"), "Not EOF error")
}
func TestDownload(t *testing.T) {
	fmt.Println("Testing torrent/download.go : Download()")
	torrentfile, err := parser.ParseFromFile("../test_torrents/ubuntu.iso.torrent")
	assert.Nil(t, err, "Opening torrent file failed.")
	u, err := url.Parse(torrentfile.Announce[0])
	statusreport := tracker.GetRandomClientReport()
	resp, err := tracker.GetPeers(u, statusreport)
	assert.Nil(t, err, "GetPeers returned error")
	assert.NotEmpty(t, resp.Peers, "Empty Peer list")
	for _, peer := range resp.Peers {
		err = Download(peer, statusreport)
		fmt.Println(err)
		assert.Nil(t, err, "Download returned error", err)
		if err==nil{
			break
		}
	}
}