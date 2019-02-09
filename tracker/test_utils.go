package tracker

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"github.com/concurrency-8/parser"
)

// Files for functions used while testing. Can be used in other packages.

// GetMockConnectResponseBuf returns a test buffer with the input transactionID and connectionID
func GetMockConnectResponseBuf(transactionID uint32, connectionID uint64) bytes.Buffer {
	var mockConnectResponseBuf bytes.Buffer

	binary.Write(&mockConnectResponseBuf, binary.BigEndian, uint32(0)) // action=0 for connect response
	binary.Write(&mockConnectResponseBuf, binary.BigEndian, transactionID)
	binary.Write(&mockConnectResponseBuf, binary.BigEndian, connectionID)

	return mockConnectResponseBuf
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

// GetRandomClientReport gives a test ClientStatusReport object pointer.
func GetRandomClientReport() (report *ClientStatusReport) {

	torrent, _ := parser.ParseFromFile(GetTorrentFileList()[0])
	report = &ClientStatusReport{}
	report.TorrentFile = torrent
	report.PeerID = string(getRandomByteArr(20))
	report.Left = torrent.Length
	report.Port = uint16(6464)
	report.Event = ""
	return report
}

// GetTorrentFileList gives the list of torrentfiles (these should be alive, in test_torrents)
func GetTorrentFileList() []string {
	return []string{"../test_torrents/ubuntu.iso.torrent", "../test_torrents/big-buck-bunny.torrent"}
}
