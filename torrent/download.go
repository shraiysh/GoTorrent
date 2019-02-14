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
	onWholeMessage(conn, handler)
}

func msgHandler(msg []byte){
	

}

// onWholeMessage sends complete messages to callback function
func onWholeMessage(conn net.Conn, msgHandler handler) error {
	buffer := new(bytes.Buffer)
	handshake := true
	resp := make([]byte, 100)

	for {
		respLen, err := conn.Read(resp)

		if err != nil {
			conn.Close() // TODO maybe a better implementation
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
			msgHandler((buffer.Bytes())[:msgLen])
			buffer = bytes.NewBuffer((buffer.Bytes())[msgLen:])
			handshake = false
		}
	}
}
