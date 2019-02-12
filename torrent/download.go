package torrent

import (
	"bytes"
	"encoding/binary"
	"net"
)

type handler func([]byte)

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
