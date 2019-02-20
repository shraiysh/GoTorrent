package torrent

import (
	"fmt"
	"github.com/concurrency-8/tracker"
	"github.com/stretchr/testify/assert"
	"github.com/concurrency-8/piece"
	"github.com/concurrency-8/queue"
	"encoding/binary"
	"math/rand"
	"net"
	"bytes"
	"testing"
)

// TestOnWholeMessage tests torrent/download.go : onWholeMessage(*kwargs)
func TestOnWholeMessage(t *testing.T) {
	fmt.Println("Testing torrent/download.go : onWholeMessage(*kwargs)")
	num := 3
	messages := make([][]byte , num)
	wholeMessage := new(bytes.Buffer)
	for i:=0 ; i < int(num) ;i++ {
		rand.Seed(int64(i))
		length := rand.Intn(256)
		var message []byte
		if i == 0 {
			message = make([]byte, length+49)
		}else {
			message = make([]byte, length+4)
		}

		rand.Read(message)
		message[0] = uint8(length)
		binary.Write(wholeMessage ,binary.BigEndian ,message) 
		messages[i] = message
	}

	client, server := net.Pipe() // create a client and server connection

	go func() {
		server.Write(wholeMessage.Bytes())
		server.Close() // close after writing out all data
	}()

	i :=0 
	err := onWholeMessage(client, func(b []byte, 
										client net.Conn , 
										pieces *piece.PieceTracker , 
										queue *queue.Queue , 
										report *tracker.ClientStatusReport) error { // mock Message Handler
		if i==0 {
			assert.Equal(t, len(b), int(messages[i][0])+49, "length not equal")
		}else {
			assert.Equal(t, len(b), int(messages[i][0])+4, "length not equal")
		}

		assert.Equal(t, b, messages[i], "message received not same")
		i++
		return assert.AnError
	} , nil , nil , nil) // TODO add tests for other message handlers

	//exclude EOF errors, due to closing a connection.
	assert.Equal(t, err, fmt.Errorf("EOF"), "Not EOF error")
	assert.Equal(t, i, int(num), "Number of messages received are not equal")
}

/*
func TestDownload(t *testing.T) {
	fmt.Println("Testing torrent/download.go : Download()")
	torrentfile, err := parser.ParseFromFile("../test_torrents/ubuntu.iso.torrent")
	assert.Nil(t, err, "Opening torrent file failed.")
	u, err := url.Parse(torrentfile.Announce[0])
	assert.Nil(t, err, "Parsing announce URL failed")
	statusreport := tracker.GetRandomClientReport()
	resp, err := tracker.GetPeers(u, statusreport)
	assert.Nil(t, err, "GetPeers returned error")
	assert.NotEmpty(t, resp.Peers, "Empty Peer list")
	for _, peer := range resp.Peers {
		err = Download(peer, statusreport)
		fmt.Println(err)
		assert.Nil(t, err, "Download returned error", err)
		if err == nil {
			break
		}
	}
}
*/

func TestChokeHandler(t *testing.T){
	client, _ := net.Pipe()

	ChokeHandler(client)

	_ , err :=  client.Read(make([]byte,4))
	assert.Equal(t, err, fmt.Errorf("io: read/write on closed pipe"), "ChokeHandler failed")
}

// func TestmsgHandler(t *testing.T){

// 	chokeMessage , err := BuildChoke() 
// }