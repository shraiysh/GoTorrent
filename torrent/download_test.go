package torrent

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"sync"
	"testing"

	"github.com/concurrency-8/parser"
	"github.com/concurrency-8/piece"
	"github.com/concurrency-8/queue"
	"github.com/concurrency-8/tracker"
	"github.com/stretchr/testify/assert"
)

func getLog() Log {
	Info := log.New(os.Stdout, "Testing ", 0)
	Error := log.New(os.Stderr, "Testing ", 0)
	return Log{
		Info:  Info,
		Error: Error,
	}
}

/*
// TestOnWholeMessage tests torrent/download.go : onWholeMessage(*kwargs)
func TestOnWholeMessage(t *testing.T) {
	fmt.Println("Testing torrent/download.go : onWholeMessage(*kwargs)")
	num := 20
	messages := make([][]byte, num)
	wholeMessage := new(bytes.Buffer)
	for i := 0; i < int(num); i++ {
		rand.Seed(int64(i))
		length := rand.Intn(256)
		var message []byte
		if i == 0 {
			message = make([]byte, length+49)
		} else {
			message = make([]byte, length+4)
		}

		rand.Read(message)
		message[0] = uint8(length)
		binary.Write(wholeMessage, binary.BigEndian, message)
		messages[i] = message
	}

	client, server := net.Pipe() // create a client and server connection

	go func() {
		server.Write(wholeMessage.Bytes())
		server.Close() // close after writing out all data
	}()

	i := 0
	_, err := onWholeMessage(tracker.Peer{}, client, func(peer tracker.Peer, b []byte,
		client net.Conn,
		pieces *piece.PieceTracker,
		queue *queue.Queue,
		report *tracker.ClientStatusReport) error { // mock Message Handler
		if i == 0 {
			assert.Equal(t, len(b), int(messages[i][0])+49, "length not equal")
		} else {
			assert.Equal(t, len(b), int(messages[i][0])+4, "length not equal")
		}

		assert.Equal(t, b, messages[i], "message received not same")
		i++
		return assert.AnError
	}, piece.NewPieceTracker(tracker.GetRandomClientReport().TorrentFile), nil, nil) // TODO add tests for other message handlers

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

// TestChokeHandler tests handling choking protocol - TODO : Adapt to the new function
// The new function tries to handshake again, just in case the peer decides to unchoke
// func TestChokeHandler(t *testing.T) {
// 	client, _ := net.Pipe()

// 	ChokeHandler(tracker.Peer{}, client, nil, nil)

// 	_, err := client.Read(make([]byte, 4))
// 	assert.Equal(t, err, fmt.Errorf("io: read/write on closed pipe"), "ChokeHandler failed")
// }

// TestUnchokeHandler tests handling of unchoking protocol
func TestUnChokeHandler(t *testing.T) {
	queue := queue.NewQueue(parser.TorrentFile{})
	UnchokeHandler(tracker.Peer{}, nil, nil, queue, getLog())
	assert.Equal(t, queue.Choked, false, "Choked attribute not set properly")
}

//TestRequestPiece tests function RequestPiece
func TestRequestPiece(t *testing.T) {

	file, _ := parser.ParseFromFile(parser.GetTorrentFileList()[0])
	pieces := piece.NewPieceTracker(file)
	queue := queue.NewQueue(file)
	queue.Choked = false
	pieceBlock := parser.RandomPieceBlock(file)
	queue.Enqueue(pieceBlock.Index)
	length := queue.Length()
	client, server := net.Pipe()
	fmt.Println(pieceBlock)
	go func() {
		for i := 0; i < length; i++ {
			RequestPiece(tracker.Peer{}, server, pieces, queue, getLog())
		}
		defer server.Close()
	}()

	for i := 0; i < length; i++ {
		resp := make([]byte, 17)
		respLen, _ := client.Read(resp)
		assert.Equal(t, respLen, 17, "Full message not received")
		size, id, payload := ParseMsg(bytes.NewBuffer(resp))
		assert.Equal(t, size, uint32(13), "Request: Size not equal")
		assert.Equal(t, id, uint8(6), "Request: Message ID different")
		assert.Equal(t, uint32(payload["index"].(uint32)), pieceBlock.Index, "Request: index field of payload not same")
		assert.Equal(t, uint32(payload["begin"].(uint32)), uint32(i)*parser.BLOCK_LEN, "Request: begin field of payload not same")
	}
}

func TestHaveHandler(t *testing.T) {
	var flag sync.WaitGroup
	flag.Add(1)
	fmt.Println("Testing torrent/download.go : HaveHandler")
	file, _ := parser.ParseFromFile(parser.GetTorrentFileList()[0])
	pieces := piece.NewPieceTracker(file)
	queue := queue.NewQueue(file)
	queue.Choked = false
	pieceBlock := parser.RandomPieceBlock(file)
	client, server := net.Pipe()
	actualsamplemsg, err := BuildHave(pieceBlock.Index)
	assert.Nil(t, err, "error writing to Buffer in BuildHave")
	go func() {
		resp := make([]byte, 20)
		_, err = server.Write(actualsamplemsg.Bytes())
		flag.Wait()
		respLen, err := server.Read(resp)
		assert.Nil(t, err, "Error reading from server")
		size, id, _ := ParseMsg(bytes.NewBuffer(resp[:respLen]))
		assert.Equal(t, uint8(6), id, "Invalid id after reading from pipe.")
		assert.Equal(t, uint32(13), size, "Invalid size")
		defer server.Close()

	}()
	resp := make([]byte, 20)
	buffer := new(bytes.Buffer)
	respLen, err := client.Read(resp)
	flag.Done()
	assert.Nil(t, err, "Error reading from pipe")
	err = binary.Write(buffer, binary.BigEndian, resp[:respLen])
	assert.Nil(t, err, "Error writing to buffer.")
	size, id, payload := ParseMsg(buffer)
	assert.Equal(t, uint8(4), id, "Invalid id after reading from pipe.")
	assert.Equal(t, uint32(5), size, "Invalid size")
	assert.NotEmpty(t, payload["payload"], "Empty piece index in payload.")
	var pieceIndex uint32
	pieceIndex, err = HaveHandler(tracker.Peer{}, client, pieces, queue, payload, getLog())
	assert.Nil(t, err, "Error in HaveHandler")
	assert.Equal(t, pieceBlock.Index, pieceIndex, "Piece Index doesn't match.")
	assert.True(t, pieces.Requested[pieceBlock.Index][0], "Requested not set.")
}

func TestBitFieldHandler(t *testing.T) {
	var flag sync.WaitGroup
	flag.Add(1)
	fmt.Println("Testing torrent/download.go : BitFieldHandler")
	file, _ := parser.ParseFromFile(parser.GetTorrentFileList()[0])
	//each piece has 20 byte hash, irrespective of size.
	npieces := uint32(len(file.Piece) / 20)
	nbytes := uint(math.Ceil(float64(npieces) / float64(8)))
	msg := new(bytes.Buffer)
	binary.Write(msg, binary.BigEndian, uint32(nbytes+1))
	binary.Write(msg, binary.BigEndian, uint8(5))
	binary.Write(msg, binary.BigEndian, getRandomByteArr(nbytes))
	actualmsg := msg.Bytes()
	pieces := piece.NewPieceTracker(file)
	queue := queue.NewQueue(file)
	queue.Choked = false
	client, server := net.Pipe()
	go func() {
		resp := make([]byte, nbytes+1)
		_, err := server.Write(actualmsg)
		flag.Wait()
		assert.Nil(t, err, "Error writing to pipe.")
		respLen, err := server.Read(resp)
		assert.Nil(t, err, "Error reading from server")
		size, id, _ := ParseMsg(bytes.NewBuffer(resp[:respLen]))
		assert.Equal(t, uint8(6), id, "Invalid id after reading from pipe.")
		assert.Equal(t, uint32(13), size, "Invalid size")
		defer server.Close()

	}()
	resp := make([]byte, nbytes+10)
	respLen, err := client.Read(resp)
	flag.Done()
	assert.Nil(t, err, "Error reading from Pipe")
	buffer := new(bytes.Buffer)
	err = binary.Write(buffer, binary.BigEndian, int32(nbytes+1))
	assert.Nil(t, err, "Error writing to buffer.")
	err = binary.Write(buffer, binary.BigEndian, int8(5))
	assert.Nil(t, err, "Error writing to buffer.")
	err = binary.Write(buffer, binary.BigEndian, resp[:respLen])
	assert.Nil(t, err, "Error writing to buffer.")
	size, id, payload := ParseMsg(buffer)
	assert.Equal(t, uint8(5), id, "Invalid id after reading from Pipe")
	assert.Equal(t, uint32(nbytes+1), size, "Invalid size")
	assert.NotEmpty(t, payload["payload"], "Empty pieces in payload")
	err = BitFieldHandler(tracker.Peer{}, client, pieces, queue, payload, getLog())
	assert.Nil(t, err, "Error in BitFieldHandler")
	// For each item in the queue, assert into the received field.
	for i := 0; queue.Length() > 0; i++ {
		nextitem, err := queue.Peek()
		assert.Nil(t, err, "Error peeking into queue")
		err = queue.Dequeue()
		assert.Nil(t, err, "Error Dequeueing from queue")
		index := nextitem.Index / 8
		//index into the byte
		offset := nextitem.Index % 8
		//offset in the byte
		p := uint8(1 << (7 - offset))
		//p is used to perform bitwise and on the byte value.
		assert.Equal(t, p, uint8(actualmsg[index])&p)
	}

}
