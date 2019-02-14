package torrent

import (
	"bytes"
	"encoding/binary"
	"net"
	"fmt"
	//"github.com/concurrency-8/parser"
	"github.com/concurrency-8/tracker"
	//"io/ioutil"
)

type handler func([]byte)

// MakeHandshake is a function that handshakes with a peer specified by peer object.
// Concurrently call this function to establish parallel connections to many peers.
func MakeHandshake(peer tracker.Peer, report *tracker.ClientStatusReport){
	buffer, err := BuildHandshake(*report)
	if err!=nil{
		return
	}
    peerip := make([]byte, 4)
    binary.BigEndian.PutUint32(peerip, peer.IPAdress)
	service := net.TCPAddr{
		IP: peerip,
		Port: int(peer.Port),
		Zone: "",
	}

	fmt.Println(service)
	//check if peer supports tcp first. anyways, if it doesn't
	// below expression will raise an error.
	conn, err := net.DialTCP("tcp", nil, &service)
	if err!=nil{
		return
	}
	//write the handshake content into the connection.
	conn.Write(buffer.Bytes())
	//read from the connection.

	resp:= make([]byte,68)
	conn.Read(resp)
	//Just to print and check stuff.
	fmt.Println(resp)
}

// onWholeMessage sends complete messages to callback function
func onWholeMessage(conn *net.IPConn, msgHandler handler) { // TODO add an extra argument for callback function i.e msgHandler
	buffer := new(bytes.Buffer)
	handshake := true
	resp := make([]byte, 100)

	for {

		respLen, err := conn.Read(resp)

		if err != nil {
			// TODO : close the connection and return
			conn.Close()
			return
		}

		binary.Write(buffer, binary.BigEndian, resp[:respLen])

		var msgLen int

		if handshake {

			var length uint8
			binary.Read(buffer, binary.BigEndian, &length)
			msgLen = int(length + 49)
		} else {

			var length int32
			binary.Read(buffer, binary.BigEndian, &length)
			msgLen = int(length + 4)
		}

		for len(buffer.Bytes()) >= 4 && len(buffer.Bytes()) >= msgLen {
			// TODO implement msgHandler
			msgHandler((buffer.Bytes())[:msgLen])
			buffer = bytes.NewBuffer((buffer.Bytes())[msgLen:])
			handshake = false
		}
	}

}
