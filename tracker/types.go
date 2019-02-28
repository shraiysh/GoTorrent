package tracker

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"github.com/concurrency-8/parser"
)

// ConnectResponse is struture to hoild details from ConnectResponse
type ConnectResponse struct {
	Action        uint32
	TransactionID uint32
	ConnectionID  uint64
}

// GetMockConnectResponseBuf returns a test buffer with the input transactionID and connectionID
func GetMockConnectResponseBuf(transactionID uint32, connectionID uint64) bytes.Buffer {
	var mockConnectResponseBuf bytes.Buffer

	binary.Write(&mockConnectResponseBuf, binary.BigEndian, uint32(0)) // action=0 for connect response
	binary.Write(&mockConnectResponseBuf, binary.BigEndian, transactionID)
	binary.Write(&mockConnectResponseBuf, binary.BigEndian, connectionID)

	return mockConnectResponseBuf
}

// AnnounceResponse is structure to hold details from announce request sent to tracker
type AnnounceResponse struct {
	Action        uint32
	TransactionID uint32
	Leechers      uint32
	Seeders       uint32
	Complete      uint   `bencode:"complete"`
	Downloaded    uint   `bencode:"downloaded"`
	Incomplete    uint   `bencode:"incomplete"`
	Interval      uint32 `bencode:"interval"`
	MinInterval   uint   `bencode:"min interval"`
	PeerBytes     []byte `bencode:"peers"`
	Peers         []Peer
}

// GetMockAnnounceResponseBuf returns a test buffer with input transactionID, interval, leechers and seeders
func GetMockAnnounceResponseBuf(transactionID, interval, leechers, seeders uint32, peers []Peer) bytes.Buffer {
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

// Peer is a structure contains IP Address of a peer
type Peer struct {
	IPAdress uint32
	Port     uint16
}

// ClientStatusReport is a structure storing current status for client and relevant information
type ClientStatusReport struct {
	Event       string
	TorrentFile parser.TorrentFile
	PeerID      string
	Port        uint16
	Uploaded    uint64
	Downloaded  uint64
	Left        uint64
	Data        []parser.Piece // This is for seeding
}

// GetRandomClientReport gives a test ClientStatusReport object pointer.
func GetRandomClientReport() (report *ClientStatusReport) {

	torrent, _ := parser.ParseFromFile(parser.GetTorrentFileList()[0])
	report = &ClientStatusReport{}
	report.TorrentFile = torrent
	report.PeerID = string(getRandomByteArr(20))
	report.Left = torrent.Length
	report.Port = uint16(6464)
	report.Event = ""
	return report
}
