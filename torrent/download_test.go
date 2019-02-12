package torrent

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net"
	"testing"
)

// TestOnWholeMessage tests torrent/download.go : onWholeMessage(*kwargs)
func TestOnWholeMessage(t *testing.T) {
	fmt.Println("Testing torrent/download.go : onWholeMessage(*kwargs)")

	length := rand.Intn(100) // generate random handshake message
	message := make([]byte, length+49)
	rand.Read(message)
	message[0] = uint8(length)

	client, server := net.Pipe() // create a client and server connection

	go func() {
		server.Write(message)
		server.Close() // close after writing out all data
	}()

	onWholeMessage(client, func(b []byte) { // mock Message Handler
		assert.Equal(t, len(b), int(b[0]+49), "length not equal")
		assert.Equal(t, b, message, "message received not same")
		fmt.Println("PASS")
	})

}
