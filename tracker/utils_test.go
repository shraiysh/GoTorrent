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

func getMockConnectResponseBuf(transactionID uint32, connectionID uint64) bytes.Buffer {
	var mockConnectResponseBuf bytes.Buffer
	writer := bufio.NewWriter(&mockConnectResponseBuf)

	binary.Write(writer, binary.BigEndian, uint32(0)) // action=0 for connect response
	binary.Write(writer, binary.BigEndian, transactionID)
	binary.Write(writer, binary.BigEndian, connectionID)
	writer.Flush()

	return mockConnectResponseBuf
}

func getMockAnnounceResponseBuf(transactionID, interval, leechers, seeders uint32, peers map[uint32]uint16) bytes.Buffer {
	var mockAnnounceResponseBuf bytes.Buffer
	writer := bufio.NewWriter(&mockAnnounceResponseBuf)

	binary.Write(writer, binary.BigEndian, uint32(1)) // action=1 for announce response
	binary.Write(writer, binary.BigEndian, transactionID)
	binary.Write(writer, binary.BigEndian, interval)
	binary.Write(writer, binary.BigEndian, leechers)
	binary.Write(writer, binary.BigEndian, seeders)
	for ip, port := range peers {
		binary.Write(writer, binary.BigEndian, ip)
		binary.Write(writer, binary.BigEndian, port)
	}

	writer.Flush()
	return mockAnnounceResponseBuf
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
	writer := bufio.NewWriter(&mockResponseBuf)

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
	trID := rand.Uint32()
	connID := rand.Uint64()
	mockConnRespBuf := getMockConnectResponseBuf(trID, connID)

	mockConnResp := ParseConnResp(mockConnRespBuf) // Object of type ConnectResponse

	assert.Equal(t, mockConnResp.action, uint32(0), "Action for connect response must be uint32(0)")
	assert.Equal(t, mockConnResp.transactionID, trID, "Unable to detect transactionID in connection response")
	assert.Equal(t, mockConnResp.connectionID, connID, "Unable to detect connectionID in connection response")
	fmt.Println("PASS")
}

func getRandomTorrent() parser.TorrentFile {
	return parser.TorrentFile{}
}

func TestBuildAnnounceReq(t *testing.T) {
	fmt.Print("Testing tracker/utils.go : BuildAnnounceReq(): ")

	connID := rand.Uint64()
	var torrent parser.TorrentFile = getRandomTorrent()
	port := uint16(6464)

	announceReqBuf := BuildAnnounceReq(connID, torrent, port)
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

	// connectionID
	binary.Read(announceReqReader, binary.BigEndian, &tempUint64)
	assert.Equal(t, connID, tempUint64, errorMsg("connectionID"))

	// action
	binary.Read(announceReqReader, binary.BigEndian, &tempUint32)
	assert.Equal(t, uint32(1), tempUint32, errorMsg("action"))

	// Cannot check for transactionID
	binary.Read(announceReqReader, binary.BigEndian, &tempUint32)

	// InfoHash
	binary.Read(announceReqReader, binary.BigEndian, &temp20ByteArr)
	var infoHash [20]byte
	copy(infoHash[:], torrent.InfoHash)
	assert.Equal(t, infoHash, temp20ByteArr, errorMsg("torrent.InfoHash"))

	// Cannot check for peerID
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
	fmt.Print("Testing tracker/utils.go : ParseAnnounceResp() : ")
	transactionID, interval, leechers, seeders := rand.Uint32(), rand.Uint32(), rand.Uint32(), rand.Uint32()
	peers := make(map[uint32]uint16)
	for i := 0; i < rand.Intn(10000); i++ {
		peers[rand.Uint32()] = uint16(rand.Intn(9000) + 1000)
	}

	mockAnnounceResponseBuf := getMockAnnounceResponseBuf(transactionID, interval, leechers, seeders, peers)

	announceResponse := ParseAnnounceResp(mockAnnounceResponseBuf)

	getErrorMsg := func(err string) string {
		return err + ": Error in ParseAnnounceResp()"
	}

	// Checking for parsed parameters
	assert.Equal(t, transactionID, announceResponse.transactionID, getErrorMsg("transactionID"))
	assert.Equal(t, interval, announceResponse.interval, getErrorMsg("interval"))
	assert.Equal(t, leechers, announceResponse.leechers, getErrorMsg("leechers"))
	assert.Equal(t, seeders, announceResponse.seeders, getErrorMsg("seeders"))

	// Checking: every peer created is parsed
	for k, v := range peers {
		assert.Contains(t, announceResponse.peers, k, getErrorMsg("keyNotParsedInPeers"))
		assert.Equal(t, v, announceResponse.peers[k], getErrorMsg("ValueMismatchInPeers"))
	}
	// Checking: No extra peer is parsed
	for k := range announceResponse.peers {
		assert.Contains(t, peers, k, getErrorMsg("ExtraKeyParsed"))
	}

	fmt.Println("PASS")
}
