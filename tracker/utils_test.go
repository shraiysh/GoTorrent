// These are tests for the functions in the file tracker/utils.go
// Run these tests with `go test` in the package directory

package tracker

import (
	"./../parser"
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"math/rand"
	"testing"
)

func getMockConnectResponseBuf(transactionId uint32, connectionId uint64) bytes.Buffer {
	var mockResponseBuf bytes.Buffer
	var writer *bufio.Writer = bufio.NewWriter(&mockResponseBuf)

	binary.Write(writer, binary.BigEndian, uint32(0)) // action=0 for connect response
	binary.Write(writer, binary.BigEndian, transactionId)
	binary.Write(writer, binary.BigEndian, connectionId)
	writer.Flush()

	return mockResponseBuf
}

func TestBuildConnReq(t *testing.T) {
	fmt.Print("Testing tracker/utils.go : BuildConnReq(): ")
	req := BuildConnReq()
	errorMessage := "Invalid Connection Request for tracker"
	assert.Equal(t, req, []byte{0x00, 0x00, 0x04, 0x17, 0x27, 0x10, 0x19, 0x80, 0x00, 0x00, 0x00, 0x00, 0xa6, 0xec, 0x6b, 0x7d}, errorMessage)
	assert.NotEqual(t, req, []byte{0x01, 0x00, 0x04, 0x17, 0x27, 0x10, 0x19, 0x80, 0x00, 0x00, 0x00, 0x00, 0xa6, 0xec, 0x6b, 0x7d}, errorMessage)
	fmt.Println("PASS")
}

func TestRespType(t *testing.T) {
	fmt.Print("Testing tracker/utils.go : RespType(): ")
	var mockResponseBuf bytes.Buffer
	var writer *bufio.Writer = bufio.NewWriter(&mockResponseBuf)

	// Mock response contains only action - announce
	binary.Write(writer, binary.BigEndian, uint32(1))
	writer.Flush()
	assert.Equal(t, RespType(mockResponseBuf), "announce", "Unable to detect \"announce\" response when action=1")

	// Mock connect response
	mockResponseBuf = getMockConnectResponseBuf(rand.Uint32(), rand.Uint64())
	assert.Equal(t, RespType(mockResponseBuf), "connect", "Unable to detect \"connect\" response when action=0")

	mockResponseBuf.Reset()

	// Mock response has 16 bytes, first 4 bytes show action - announce
	binary.Write(writer, binary.BigEndian, uint32(1)) // 4 bytes written
	for i := 0; i < 3; i++ {                          // Next 12 bytes = 3 uint32
		binary.Write(writer, binary.BigEndian, uint32(rand.Uint32()))
	}
	writer.Flush()
	assert.Equal(t, RespType(mockResponseBuf), "announce", "Unable to detect \"announce\" response when action=1")

	fmt.Println("PASS")
}

func TestParseConnResp(t *testing.T) {
	fmt.Print("Testing tracker/utils.go : ParseConnResp(): ")
	var trId uint32 = rand.Uint32()
	var connId uint64 = rand.Uint64()
	var mockConnRespBuf bytes.Buffer = getMockConnectResponseBuf(trId, connId)

	var mockConnResp ConnectResponse = ParseConnResp(mockConnRespBuf)

	assert.Equal(t, mockConnResp.action, uint32(0), "Action for connect response must be uint32(0)")
	assert.Equal(t, mockConnResp.transactionId, trId, "Unable to detect transactionId in connection response")
	assert.Equal(t, mockConnResp.connectionId, connId, "Unable to detect connectionId in connection response")
	fmt.Println("PASS")
}

func getRandomTorrent() parser.TorrentFile {
	return parser.TorrentFile{}
}

func TestBuildAnnounceReq(t *testing.T) {
	fmt.Print("Testing tracker/utils.go : BuildAnnounceReq(): ")

	var connId uint64 = rand.Uint64()
	var torrent parser.TorrentFile = getRandomTorrent()
	var port uint16 = uint16(6464)

	var announceReqBuf bytes.Buffer = BuildAnnounceReq(connId, torrent, port)
	var announceReqReader io.Reader = bytes.NewReader(announceReqBuf.Bytes())

	// Temporary variables to store data read from generated buffer
	var tempUint64 uint64
	var tempUint32 uint32
	var temp20ByteArr [20]byte
	var tempInt32 int32
	var tempUint16 uint16

	errorMsg := func(varName string) string {
		return varName + ": not set properly in Announce Request"
	}

	// connectionId
	binary.Read(announceReqReader, binary.BigEndian, &tempUint64)
	assert.Equal(t, connId, tempUint64, errorMsg("connectionId"))

	// action
	binary.Read(announceReqReader, binary.BigEndian, &tempUint32)
	assert.Equal(t, uint32(1), tempUint32, errorMsg("action"))

	// Cannot check for transactionId
	binary.Read(announceReqReader, binary.BigEndian, &tempUint32)

	// InfoHash
	binary.Read(announceReqReader, binary.BigEndian, &temp20ByteArr)
	var infoHash [20]byte
	copy(infoHash[:], torrent.InfoHash)
	assert.Equal(t, infoHash, temp20ByteArr, errorMsg("torrent.InfoHash"))

	// Cannot check for peerId
	binary.Read(announceReqReader, binary.BigEndian, &temp20ByteArr)

	// downloaded
	binary.Read(announceReqReader, binary.BigEndian, &tempUint64)
	assert.Equal(t, uint64(0), tempUint64, errorMsg("downloaded"))

	// left
	binary.Read(announceReqReader, binary.BigEndian, &tempUint64)
	assert.Equal(t, torrent.Length, tempUint64, errorMsg("left"))

	// uploaded
	binary.Read(announceReqReader, binary.BigEndian, &tempUint64)
	assert.Equal(t, uint64(0), tempUint64, errorMsg("uploaded"))

	// event
	binary.Read(announceReqReader, binary.BigEndian, &tempUint32)
	assert.Equal(t, uint32(0), tempUint32, errorMsg("event"))

	// Ip address
	binary.Read(announceReqReader, binary.BigEndian, &tempUint32)
	assert.Equal(t, uint32(0), tempUint32, errorMsg("Ip address"))

	// Cannot check key
	binary.Read(announceReqReader, binary.BigEndian, &tempUint32)

	// num want
	binary.Read(announceReqReader, binary.BigEndian, &tempInt32)
	assert.Equal(t, int32(-1), tempInt32, errorMsg("num_want"))

	// port
	binary.Read(announceReqReader, binary.BigEndian, &tempUint16)
	assert.Equal(t, port, tempUint16, errorMsg("port"))

	fmt.Println("PASS")
}

func TestParseAnnounceResp(t *testing.T) {
	// TODO Implement this test for tracker/utils.go: ParseAnnounceResp(bytes.Buffer)
}
