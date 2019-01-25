package tracker

import (
	"../parser"
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"os"
)

// BuildConnReq is the first connection request for tracker
func BuildConnReq() []byte {
	buffer := &bytes.Buffer{}
	buffer.Write([]byte{0x00, 0x00, 0x04, 0x17,
		0x27, 0x10, 0x19, 0x80,
		0x00, 0x00, 0x00, 0x00,
		0xa6, 0xec, 0x6b, 0x7d})

	return buffer.Bytes()
}

// RespType is for decoding the responses received from socket
func RespType(response bytes.Buffer) string {
	action := binary.BigEndian.Uint32(response.Bytes()[0:4])
	if action == 0 {
		return "connect"
	}
	return "announce"
}

// ParseConnResp parses the connection request and returns action, transactionId and connectionId
func ParseConnResp(response bytes.Buffer) (uint32, uint32, uint64) {
	responseBytes := response.Bytes()
	action := binary.BigEndian.Uint32(responseBytes[0:4])
	transactionId := binary.BigEndian.Uint32(responseBytes[4:8])
	connectionId := binary.BigEndian.Uint64(responseBytes[8:])
	return action, transactionId, connectionId
}

func writeUint64ToBuffer(buf *bytes.Buffer, value uint64) {
	temp := make([]byte, 8)
	binary.BigEndian.PutUint64(temp, value)
	buf.Write(temp)
}

func getRandomByteArr(size uint) []byte {
	temp := make([]byte, size)
	_, err := rand.Read(temp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to generate Crypto random byte array")
	}
	return temp
}

// BuildAnnounceReq builds an announce request where we tell the tracker which files we're interested in
func BuildAnnounceReq(connectionId uint64, torrent parser.TorrentFile, port uint16) bytes.Buffer {
	buffer := bytes.NewBuffer(make([]byte, 0, 98))

	// connection id
	writeUint64ToBuffer(buffer, connectionId)

	// action
	buffer.Write([]byte{1, 0, 0, 0})

	// transaction id
	buffer.Write(getRandomByteArr(4))

	// info hash
	buffer.WriteString(torrent.InfoHash)

	// peer id
	buffer.Write(getRandomByteArr(20))

	// downloaded
	buffer.Write(make([]byte, 0, 8))

	// left
	writeUint64ToBuffer(buffer, torrent.Length)

	// uploaded
	buffer.Write(make([]byte, 0, 8))

	// event
	buffer.Write(make([]byte, 0, 4))

	// ip address
	buffer.Write(make([]byte, 0, 4))

	// key
	buffer.Write(getRandomByteArr(4))

	// num want
	numWant := make([]byte, 4)
	binary.PutVarint(numWant, -1)
	buffer.Write(numWant)

	// port
	portArr := make([]byte, 2)
	binary.BigEndian.PutUint16(portArr, port)
	buffer.Write(portArr)

	return *buffer
}
