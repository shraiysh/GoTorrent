// These are tests for the functions in the file tracker/utils.go
// Run these tests with `go test` in the package directory

package tracker

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/concurrency-8/parser"
	"github.com/stretchr/testify/assert"
	"io"
	"math/rand"
	"net/url"
	"testing"
)

func getTorrentFileList() []string {
	return []string{"../test_torrents/ubuntu.iso.torrent", "../test_torrents/big-buck-bunny.torrent"}
}

func getErrorMsg(varName, functionName string) string {
	return varName + ": not set properly in " + functionName + ". Tip: You might want to check if your network allows torrenting!"
}

func getMockConnectResponseBuf(transactionID uint32, connectionID uint64) bytes.Buffer {
	var mockConnectResponseBuf bytes.Buffer
	writer := bufio.NewWriter(&mockConnectResponseBuf)

	binary.Write(writer, binary.BigEndian, uint32(0)) // action=0 for connect response
	binary.Write(writer, binary.BigEndian, transactionID)
	binary.Write(writer, binary.BigEndian, connectionID)
	writer.Flush()

	return mockConnectResponseBuf
}

func getMockAnnounceResponseBuf(transactionID, interval, leechers, seeders uint32, peers []Peer) bytes.Buffer {
	var mockAnnounceResponseBuf bytes.Buffer
	writer := bufio.NewWriter(&mockAnnounceResponseBuf)

	binary.Write(writer, binary.BigEndian, uint32(1)) // action=1 for announce response
	binary.Write(writer, binary.BigEndian, transactionID)
	binary.Write(writer, binary.BigEndian, interval)
	binary.Write(writer, binary.BigEndian, leechers)
	binary.Write(writer, binary.BigEndian, seeders)

	for i := 0; i < len(peers); i++ {
		binary.Write(writer, binary.BigEndian, peers[i].IPAdress)
		binary.Write(writer, binary.BigEndian, peers[i].Port)
	}

	writer.Flush()
	return mockAnnounceResponseBuf
}

func TestBuildConnReq(t *testing.T) {
	fmt.Print("Testing tracker/utils.go : buildConnReq(): ")
	req := buildConnReq()
	errorMessage := "Invalid Connection Request for tracker"
	assert.Equal(t, req[:12], []byte{0x00, 0x00, 0x04, 0x17, 0x27, 0x10, 0x19, 0x80, 0x00, 0x00, 0x00, 0x00}, errorMessage)
	assert.NotEqual(t, req[:12], []byte{0x01, 0x00, 0x04, 0x17, 0x27, 0x10, 0x19, 0x80, 0x00, 0x00, 0x00, 0x00}, errorMessage)

	fmt.Println("PASS")
}

func TestRespType(t *testing.T) {

	fmt.Print("Testing tracker/utils.go : respType(): ")

	var mockResponseBuf bytes.Buffer
	writer := bufio.NewWriter(&mockResponseBuf)

	// Mock response contains only action - announce
	binary.Write(writer, binary.BigEndian, uint32(1))
	writer.Flush()
	assert.Equal(t, respType(mockResponseBuf), "announce", "Unable to detect \"announce\" response when action=1")

	// Mock connect response
	mockResponseBuf = getMockConnectResponseBuf(rand.Uint32(), rand.Uint64())
	assert.Equal(t, respType(mockResponseBuf), "connect", "Unable to detect \"connect\" response when action=0")

	mockResponseBuf.Reset()

	// Mock response has 16 bytes, first 4 bytes show action - announce
	binary.Write(writer, binary.BigEndian, uint32(1)) // 4 bytes written
	for i := 0; i < 3; i++ {                          // Next 12 bytes = 3 uint32
		binary.Write(writer, binary.BigEndian, uint32(rand.Uint32()))
	}
	writer.Flush()

	assert.Equal(t, respType(mockResponseBuf), "announce", "Unable to detect \"announce\" response when action=1")

	fmt.Println("PASS")
}

func TestParseConnResp(t *testing.T) {

	fmt.Print("Testing tracker/utils.go : parseConnResp(): ")
	trID := rand.Uint32()
	connID := rand.Uint64()
	mockConnRespBuf := getMockConnectResponseBuf(trID, connID)

	mockConnResp := parseConnResp(mockConnRespBuf) // Object of type ConnectResponse

	assert.Equal(t, mockConnResp.Action, uint32(0), "Action for connect response must be uint32(0)")
	assert.Equal(t, mockConnResp.TransactionID, trID, "Unable to detect transactionID in connection response")
	assert.Equal(t, mockConnResp.ConnectionID, connID, "Unable to detect connectionID in connection response")
	fmt.Println("PASS")
}

func getRandomTorrent() parser.TorrentFile {
	return parser.TorrentFile{}
}

func getRandomClientReport() (report *ClientStatusReport) {

	torrent, _ := parser.ParseFromFile(getTorrentFileList()[0])
	report = &ClientStatusReport{}
	report.TorrentFile = torrent
	report.PeerID = string(getRandomByteArr(20))
	report.Left = torrent.Length
	report.Port = uint16(6464)
	report.Event = ""
	return report
}

func TestBuildAnnounceReq(t *testing.T) {
	fmt.Print("Testing tracker/utils.go : buildAnnounceReq(): ")

	connID := rand.Uint64()
	report := getRandomClientReport()

	announceReqBuf, _ := buildAnnounceReq(connID, report)

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

	copy(infoHash[:], report.TorrentFile.InfoHash)
	assert.Equal(t, infoHash, temp20ByteArr, errorMsg("torrent.InfoHash"))

	// Cannot check for peerID
	binary.Read(announceReqReader, binary.BigEndian, &temp20ByteArr)

	// downloaded
	binary.Read(announceReqReader, binary.BigEndian, &tempUint64)
	assert.Equal(t, uint64(0), tempUint64, errorMsg("downloaded"))

	// left
	binary.Read(announceReqReader, binary.BigEndian, &tempUint64)

	assert.Equal(t, report.Left, tempUint64, errorMsg("left"))
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
	assert.Equal(t, report.Port, tempUint16, errorMsg("port"))

	fmt.Println("PASS")
}

func TestParseAnnounceResp(t *testing.T) {
	fmt.Print("Testing tracker/utils.go : parseAnnounceResp(): ")
	transactionID, interval, leechers, seeders := rand.Uint32(), rand.Uint32(), rand.Uint32(), rand.Uint32()
	length := rand.Intn(5)
	peers := make([]Peer, length)
	for i := 0; i < length; i++ {
		peers[i].IPAdress = rand.Uint32()
		peers[i].Port = uint16(rand.Intn(9000) + 1000)
	}

	mockAnnounceResponseBuf := getMockAnnounceResponseBuf(transactionID, interval, leechers, seeders, peers)

	announceResponse := parseAnnounceResp(mockAnnounceResponseBuf)

	// Checking for parsed parameters
	assert.Equal(t, transactionID, announceResponse.TransactionID, getErrorMsg("transactionID", "TestParseAnnounceResp"))
	assert.Equal(t, interval, announceResponse.Interval, getErrorMsg("interval", "TestParseAnnounceResp"))
	assert.Equal(t, leechers, announceResponse.Leechers, getErrorMsg("leechers", "TestParseAnnounceResponse"))
	assert.Equal(t, seeders, announceResponse.Seeders, getErrorMsg("seeders", "TestParseAnnounceResponse"))

	// Checking: Number of peers received is same
	assert.Equal(t, len(peers), len(announceResponse.Peers), getErrorMsg("LengthNotEqual", "TestParseAnnounceResponse"))

	// Checking: every peer created is parsed
	for i := 0; i < length; i++ {
		assert.Equal(t, peers[i].IPAdress, announceResponse.Peers[i].IPAdress, getErrorMsg("IPMismatchInPeers", "TestParseAnnounceResponse"))
		assert.Equal(t, peers[i].Port, announceResponse.Peers[i].Port, getErrorMsg("PortMismatchInPeers", "TestParseAnnounceResponse"))

	}

	fmt.Println("PASS")
}

func TestGetPeers(t *testing.T) {

	fmt.Print("Testing tracker/utils.go : GetPeers(): ")

	for _, torrentfileName := range getTorrentFileList() {
		//torrentfile := getRandomTorrent();
		torrentfile, _ := parser.ParseFromFile(torrentfileName)
		passes := false
		for _, announceUrl := range torrentfile.Announce {

			u, err := url.Parse(announceUrl)
			if err != nil {
				fmt.Println("\nWarning:", err)
				continue
			}

			clientReport := GetClientStatusReport(torrentfile, 6881)

			_, err = GetPeers(u, clientReport)
			if err != nil {
				fmt.Println("\nWarning:", err)
				continue
			}
			passes = true
			break
		}
		assert.Equal(t, true, passes, getErrorMsg("torrentfile", "TestGetPeers"))
	}
	fmt.Println("PASS")
}
