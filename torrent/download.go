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

// Download is a function that handshakes with a peer specified by peer object.
// Concurrently call this function to establish parallel connections to many peers.
func Download(peer tracker.Peer, report *tracker.ClientStatusReport){
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
	conn, err := net.DialTCP("tcp", nil, &service)
	conn.SetKeepAlive(true)
	if err!=nil{
		return
	}
	//write the handshake content into the connection.
	conn.Write(buffer.Bytes())
	//use onWholeMessage() to read safely from the conn.
	//TODO: make a handler function and the paramter here.
	//Build will fail without this.
	onWholeMessage(conn, handler)
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
