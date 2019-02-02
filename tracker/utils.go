package tracker

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/concurrency-8/parser"
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

// ParseConnResp parses the connection request and returns action, transactionID and connectionID
func ParseConnResp(response bytes.Buffer) ConnectResponse {
	var connectionResponse ConnectResponse
	responseBytes := response.Bytes()
	connectionResponse.action = binary.BigEndian.Uint32(responseBytes[0:4])
	connectionResponse.transactionID = binary.BigEndian.Uint32(responseBytes[4:8])
	connectionResponse.connectionID = binary.BigEndian.Uint64(responseBytes[8:])
	return connectionResponse
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
func BuildAnnounceReq(connectionID uint64, torrent parser.TorrentFile, port uint16) bytes.Buffer {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)

	// connection id
	binary.Write(writer, binary.BigEndian, connectionID)

	// action
	binary.Write(writer, binary.BigEndian, uint32(1)) // announce req

	// transaction id
	binary.Write(writer, binary.BigEndian, getRandomByteArr(4))

	// info hash
	var infoHash [20]byte
	copy(infoHash[:], torrent.InfoHash)
	binary.Write(writer, binary.BigEndian, infoHash)

	// peer id
	binary.Write(writer, binary.BigEndian, getRandomByteArr(20))

	// downloaded
	binary.Write(writer, binary.BigEndian, uint64(0))

	// left
	binary.Write(writer, binary.BigEndian, torrent.Length)

	// uploaded
	binary.Write(writer, binary.BigEndian, uint64(0))

	// event
	binary.Write(writer, binary.BigEndian, uint32(0))

	// ip address
	binary.Write(writer, binary.BigEndian, uint32(0))

	// key
	binary.Write(writer, binary.BigEndian, getRandomByteArr(4))

	// num want
	binary.Write(writer, binary.BigEndian, int32(-1))

	// port
	binary.Write(writer, binary.BigEndian, port)

	writer.Flush()

	return buffer
}

// ParseAnnounceResp parses necessary details from the announce response sent by tracker
func ParseAnnounceResp(response bytes.Buffer) AnnounceResponse {
	var result AnnounceResponse

	responseBytes := response.Bytes()

	result.action = binary.BigEndian.Uint32(responseBytes[0:4])
	result.transactionID = binary.BigEndian.Uint32(responseBytes[4:8])
	result.interval = binary.BigEndian.Uint32(responseBytes[8:12])
	result.leechers = binary.BigEndian.Uint32(responseBytes[12:16])
	result.seeders = binary.BigEndian.Uint32(responseBytes[16:20])

	result.peers = make(map[uint32]uint16)
	for i := 20; i+5 < len(responseBytes); i += 6 {
		result.peers[binary.BigEndian.Uint32(responseBytes[i:i+4])] = binary.BigEndian.Uint16(responseBytes[i+4 : i+6])
	}

	return result
}
