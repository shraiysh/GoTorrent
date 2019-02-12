package torrent

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net"
	"testing"
	"time"
)

// TestOnWholeMessage tests torrent/download.go : onWholeMessage(*kwargs)
func TestOnWholeMessage(t *testing.T) {
	fmt.Println("Testing torrent/download.go : onWholeMessage(*kwargs)")

	// start a localhost server as goroutine to act as a peer
	go func() {
		ln, _ := net.Listen("tcp", ":8000")
		conn, _ := ln.Accept()
		length := rand.Intn(100)
		message := make([]byte, length+49)
		rand.Read(message)
		message[0] = uint8(length)
		conn.Write(message)
	}()

	time.Sleep(time.Second)

	// make a client to connect to server
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	assert.Equal(t, err, nil, "can't dial IP 127.0.0.1:8081")

	onWholeMessage(conn, func(b []byte) { // mock Message Handler
		assert.Equal(t, len(b), int(b[0]+49), "length not equal")
		fmt.Println("PASS")
		return
	}, true)

}
