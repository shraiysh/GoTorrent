package torrent

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/concurrency-8/tracker"
	"net"
)

type handler func([]byte, net.Conn) error

// Download is a function that handshakes with a peer specified by peer object.
// Concurrently call this function to establish parallel connections to many peers.
func Download(peer tracker.Peer, report *tracker.ClientStatusReport) error {
	buffer, err := BuildHandshake(*report)
	if err != nil {
		return err
	}
	peerip := make([]byte, 4)
	binary.BigEndian.PutUint32(peerip, peer.IPAdress)
	service := net.TCPAddr{
		IP:   peerip,
		Port: int(peer.Port),
		Zone: "",
	}
	conn, err := net.Dial("tcp", service.String())
	if err != nil {
		return err
	}
	//write the handshake content into the connection.
	_, err = conn.Write(buffer.Bytes())
	if err != nil {
		return err
	}
	//safely handle reading using onWholeMessage
	onWholeMessage(conn, msgHandler)
	return err
}
func msgHandler(msg []byte, conn net.Conn) error {

	if (len(msg) == int(uint8(msg[0]))+49) && (bytes.Equal(msg[1:20], []byte("BitTorrent protocol"))) {
		message, err := BuildInterested()
		if err != nil {
			fmt.Println("Error", err.Error())
			return err
		}
		conn.Write(message.Bytes())
	}
	/* Other non-handshake functions should follow */
	return nil

}

// onWholeMessage sends complete messages to callback function
func onWholeMessage(conn net.Conn, msgHandler handler) error {
	buffer := new(bytes.Buffer)
	handshake := true
	resp := make([]byte, 100)

	for {
		respLen, err := conn.Read(resp)
		//Please look for a better connection handling in the future.
		//Maybe use defer?
		if err != nil {
			conn.Close()
			return err
		}

		binary.Write(buffer, binary.BigEndian, resp[:respLen])

		var msgLen int

		if handshake {
			length := uint8((buffer.Bytes())[0])
			msgLen = int(length) + 49
		} else {

			length := int32((buffer.Bytes())[0])
			msgLen = int(length + 4)
		}
		for len(buffer.Bytes()) >= 4 && len(buffer.Bytes()) >= msgLen {
			// TODO implement msgHandler
			msgHandler((buffer.Bytes())[:msgLen], conn)
			buffer = bytes.NewBuffer((buffer.Bytes())[msgLen:])
			handshake = false
		}
	}
}
