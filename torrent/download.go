package main

import (
	"bytes"
	"encoding/binary"
	url "net/url"
	"net"
	"fmt"
	"github.com/concurrency-8/parser"
	"github.com/concurrency-8/tracker"
	"io/ioutil"
)

type handler func([]byte)

// MakeHandshake is this.
func MakeHandshake(u *url.URL, torrent parser.TorrentFile){
	report := tracker.GetClientStatusReport(torrent, uint16(u.Port()))
	buffer, err := BuildHandshake(*report)
	if err!=nil{
		return
	}
	service := "46.182.109.197:64806"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	if err!=nil{
		return
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err!=nil{
		return
	}
	conn.Write(buffer.Bytes())
	result, err := ioutil.ReadAll(conn)
	fmt.Println(string(result))


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
