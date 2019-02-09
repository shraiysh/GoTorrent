package torrent

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/concurrency-8/parser"
	"github.com/concurrency-8/tracker"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"testing"
)

func getTorrentFileList() []string {
	return []string{"../ubuntu.iso.torrent", "../big-buck-bunny.torrent"}
}

func getRandomClientReport() (report *tracker.ClientStatusReport) {
	torrent, _ := parser.ParseFromFile(getTorrentFileList()[0])
	report = &tracker.ClientStatusReport{}
	report.TorrentFile = torrent
	report.PeerID = string(getRandomByteArr(20))
	report.Left = torrent.Length
	report.Port = uint16(6464)
	report.Event = ""
	return report
}

func getRandomByteArr(size uint) []byte {
	temp := make([]byte, size)
	_, err := rand.Read(temp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to generate Crypto random byte array")
	}
	return temp
}

func TestBuildHandshake(t *testing.T) {
	assert := assert.New(t)
	csr := getRandomClientReport()
	handshake, err := BuildHandshake(*csr)

	assert.Nil(err)

	handshakeReader := bytes.NewReader(handshake.Bytes())

	var pstrlen uint8
	assert.Nil(binary.Read(handshakeReader, binary.BigEndian, &pstrlen))

	// Checks for the length and hence confirming that the pstrlen and pstr part are in sync
	assert.Equal(int(49+pstrlen), len(handshake.Bytes()))

	pstr := make([]byte, pstrlen)
	assert.Nil(binary.Read(handshakeReader, binary.BigEndian, &pstr))

	var reserved [8]byte
	assert.Nil(binary.Read(handshakeReader, binary.BigEndian, &reserved))

	var infohashFromBuf, infohashFromFile [20]byte
	assert.Nil(binary.Read(handshakeReader, binary.BigEndian, &infohashFromBuf))
	copy(infohashFromFile[:], csr.TorrentFile.InfoHash)
	assert.Equal(infohashFromFile, infohashFromBuf)

	var peerID [20]byte
	assert.Nil(binary.Read(handshakeReader, binary.BigEndian, &peerID))
	assert.Equal(csr.PeerID, string(peerID[:]))
}

func TestBuildKeepAlive(t *testing.T) {
	alive := BuildKeepAlive()
	assert.Equal(t, *bytes.NewBuffer(make([]byte, 4)), *alive)
}

func TestBuildChoke(t *testing.T) {
	choke, err := BuildChoke()

	assert.Nil(t, err)

	// Length of choke byte array = 5
	assert.Equal(t, 5, len(choke.Bytes()))

	chokeReader := bytes.NewReader(choke.Bytes())

	// Read length
	var length uint32
	assert.Nil(t, binary.Read(chokeReader, binary.BigEndian, &length))
	assert.Equal(t, uint32(1), length)

	// Read message type (choke)
	var messageType uint8
	assert.Nil(t, binary.Read(chokeReader, binary.BigEndian, &messageType))
	assert.Equal(t, uint8(0), messageType)
}

func TestBuildUnchoke(t *testing.T) {
	unchoke, err := BuildUnchoke()

	assert.Nil(t, err)

	// Length of unchoke byte array = 5
	assert.Equal(t, 5, len(unchoke.Bytes()))

	unchokeReader := bytes.NewReader(unchoke.Bytes())

	// Read length
	var length uint32
	assert.Nil(t, binary.Read(unchokeReader, binary.BigEndian, &length))
	assert.Equal(t, uint32(1), length)

	// Read message type (unchoke)
	var messageType uint8
	assert.Nil(t, binary.Read(unchokeReader, binary.BigEndian, &messageType))
	assert.Equal(t, uint8(1), messageType)
}

func TestBuildInterested(t *testing.T) {

	assert := assert.New(t)

	interested, err := BuildInterested()

	assert.Nil(err)
	// Length of interested byte array = 5
	assert.Equal(5, len(interested.Bytes()))

	interestedReader := bytes.NewReader(interested.Bytes())

	// Read length
	var length uint32
	assert.Nil(binary.Read(interestedReader, binary.BigEndian, &length))
	assert.Equal(uint32(1), length)

	// Read message type (interested)
	var messageType uint8
	assert.Nil(binary.Read(interestedReader, binary.BigEndian, &messageType))
	assert.Equal(uint8(2), messageType)

}

func TestBuildUninterested(t *testing.T) {
	assert := assert.New(t)

	uninterested, err := BuildUninterested()

	assert.Nil(err)
	// Length of uninterested byte array = 5
	assert.Equal(5, len(uninterested.Bytes()))

	uninterestedReader := bytes.NewReader(uninterested.Bytes())

	// Read length
	var length uint32
	assert.Nil(binary.Read(uninterestedReader, binary.BigEndian, &length))
	assert.Equal(uint32(1), length)

	// Read message type (uninterested)
	var messageType uint8
	assert.Nil(binary.Read(uninterestedReader, binary.BigEndian, &messageType))
	assert.Equal(uint8(3), messageType)
}

func TestBuildHave(t *testing.T) {
	assert := assert.New(t)

	payload := uint32(rand.Uint32())
	have, err := BuildHave(payload)

	assert.Nil(err)
	// Length of have byte array = 5
	assert.Equal(9, len(have.Bytes()))

	haveReader := bytes.NewReader(have.Bytes())

	// Read length
	var length uint32
	assert.Nil(binary.Read(haveReader, binary.BigEndian, &length))
	assert.Equal(uint32(5), length)

	// Read message type (have)
	var messageType uint8
	assert.Nil(binary.Read(haveReader, binary.BigEndian, &messageType))
	assert.Equal(uint8(4), messageType)

	// Read payload
	var payloadRead uint32
	assert.Nil(binary.Read(haveReader, binary.BigEndian, &payloadRead))
	assert.Equal(payload, payloadRead)
}
