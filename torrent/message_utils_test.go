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

	// Length of choke buffer
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

	// Length of unchoke buffer
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
	// Length of interested buffer
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
	// Length of uninterested buffer
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
	// Length of have buffer
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

func TestBuildRequest(t *testing.T) {
	assert := assert.New(t)

	piece := GetRandomPiece()

	request, err := BuildRequest(piece)

	assert.Nil(err)

	// Length of request buffer
	assert.Equal(17, len(request.Bytes()))

	requestReader := bytes.NewReader(request.Bytes())

	// Read length of message
	var length uint32
	assert.Nil(binary.Read(requestReader, binary.BigEndian, &length))
	assert.Equal(uint32(13), length)

	// Read message type
	var messageType uint8
	assert.Nil(binary.Read(requestReader, binary.BigEndian, &messageType))
	assert.Equal(uint8(6), messageType)

	// Read piece index
	var pieceIndex uint32
	assert.Nil(binary.Read(requestReader, binary.BigEndian, &pieceIndex))
	assert.Equal(piece.Index, pieceIndex)

	// Read piece begin point
	var pieceBegin uint32
	assert.Nil(binary.Read(requestReader, binary.BigEndian, &pieceBegin))
	assert.Equal(piece.Begin, pieceBegin)

	// Read piece length
	var pieceLength uint32
	assert.Nil(binary.Read(requestReader, binary.BigEndian, &pieceLength))
	assert.Equal(piece.Length, pieceLength)
}

// func TestBuildPiece(t *testing.T) {
// 	assert := assert.New(t)

// 	samplePiece := GetRandomPiece()

// 	builtPiece, err := BuildPiece(samplePiece)

// 	assert.Nil(err)

// 	// Length of piece buffer
// 	assert.Equal(len(samplePiece.Block.Bytes())+13, len(builtPiece.Bytes()))

// 	pieceReader := bytes.NewReader(builtPiece.Bytes())

// 	// Read length of message
// 	var length uint32
// 	assert.Nil(binary.Read(pieceReader, binary.BigEndian, &length))
// 	assert.Equal(uint32(len(samplePiece.Block.Bytes())+9), length)

// 	// Read message type
// 	var messageType uint8
// 	assert.Nil(binary.Read(pieceReader, binary.BigEndian, &messageType))
// 	assert.Equal(uint8(7), messageType)

// 	// Read piece index
// 	var pieceIndexRead uint32
// 	assert.Nil(binary.Read(pieceReader, binary.BigEndian, &pieceIndexRead))
// 	assert.Equal(samplePiece.Index, pieceIndexRead)

// 	// Read Begin
// 	var pieceBeginRead uint32
// 	assert.Nil(binary.Read(pieceReader, binary.BigEndian, &pieceBeginRead))
// 	assert.Equal(samplePiece.Begin, pieceBeginRead)

// 	// Read Buffer
// 	blockRead := make([]byte, len(samplePiece.Block.Bytes()))
// 	assert.Nil(binary.Read(pieceReader, binary.BigEndian, &blockRead))
// 	assert.Equal(samplePiece.Block.Bytes(), blockRead)
// }

func TestBuildCancel(t *testing.T) {
	assert := assert.New(t)

	samplePiece := GetRandomPiece()

	cancel, err := BuildCancel(samplePiece)

	assert.Nil(err)

	// Length of cancel buffer
	assert.Equal(17, len(cancel.Bytes()))

	cancelReader := bytes.NewReader(cancel.Bytes())

	// Read length of message
	var length uint32
	assert.Nil(binary.Read(cancelReader, binary.BigEndian, &length))
	assert.Equal(uint32(13), length)

	// Message type = 8
	var messageType uint8
	assert.Nil(binary.Read(cancelReader, binary.BigEndian, &messageType))
	assert.Equal(uint8(8), messageType)

	// piece index
	var pieceIndex uint32
	assert.Nil(binary.Read(cancelReader, binary.BigEndian, &pieceIndex))
	assert.Equal(samplePiece.Index, pieceIndex)

	// piece begin
	var pieceBegin uint32
	assert.Nil(binary.Read(cancelReader, binary.BigEndian, &pieceBegin))
	assert.Equal(samplePiece.Begin, pieceBegin)

	// piece length
	var pieceLength uint32
	assert.Nil(binary.Read(cancelReader, binary.BigEndian, &pieceLength))
	assert.Equal(samplePiece.Length, pieceLength)
}

func TestBuildPort(t *testing.T) {
	assert := assert.New(t)
	port := uint16(rand.Intn(90000) + 10000)

	portBuf, err := BuildPort(port)
	assert.Nil(err)

	// Length of port buffer
	assert.Equal(7, len(portBuf.Bytes()))

	portBufReader := bytes.NewReader(portBuf.Bytes())

	// Length of message
	var length uint32
	assert.Nil(binary.Read(portBufReader, binary.BigEndian, &length))
	assert.Equal(uint32(3), length)

	// Message type = 9
	var messageType uint8
	assert.Nil(binary.Read(portBufReader, binary.BigEndian, &messageType))
	assert.Equal(uint8(9), messageType)

	// port
	var portReadFromBuf uint16
	assert.Nil(binary.Read(portBufReader, binary.BigEndian, &portReadFromBuf))
	assert.Equal(port, portReadFromBuf)
}

func TestParseMsg(t *testing.T){

	// BuildChoke
	chokeMessage , _ := BuildChoke()
	size ,id , payload := ParseMsg(chokeMessage)
	assert.Equal(t , size ,int32(1), "choke : Size not equal")
	assert.Equal(t , id , int8(0), "choke : Message ID different")
	assert.Equal(t , len(payload) , 0 , "choke : length of payload not zero")

	//BuildUnchoke
	unchokeMessage , _ := BuildUnchoke()
	size ,id , payload = ParseMsg(unchokeMessage)
	assert.Equal(t , size ,int32(1), "unchoke : Size not equal")
	assert.Equal(t , id , int8(1), "unchoke : Message ID different")
	assert.Equal(t , len(payload) , 0 , "unchoke : length of payload not zero")

	// BuildInterested
	interestedMessage , _ := BuildInterested()
	size ,id , payload = ParseMsg(interestedMessage)
	assert.Equal(t , size ,int32(1), "Interested : Size not equal")
	assert.Equal(t , id , int8(2), "Interested : Message ID different")
	assert.Equal(t , len(payload) , 0 , "Interested: length of payload not zero")

	// BuildUninerested
	uninterestedMessage , _ := BuildUninterested()
	size ,id , payload = ParseMsg(uninterestedMessage)
	assert.Equal(t , size ,int32(1), "UnInterested: Size not equal")
	assert.Equal(t , id , int8(3), "Unterested: Message ID different")
	assert.Equal(t , len(payload) , 0 , "UnInterested: length of payload not zero")

	// BuildHave
	p1 := rand.Uint32()
	haveMessage , _ := BuildHave(p1)
	size ,id , payload = ParseMsg(haveMessage)
	assert.Equal(t , size ,int32(5), "Have: Size not equal")
	assert.Equal(t , id , int8(4), "Have: Message ID different")
	var p2 uint32
	binary.Read(payload["payload"].(*bytes.Buffer), binary.BigEndian , &p2)
	assert.Equal(t , p2 , p1 , "Have: length of payload not zero")

	file , _ := parser.ParseFromFile(parser.GetTorrentFileList()[0])
	pieceBlock := parser.RandomPieceBlock(file)

	// BuildRequest
	requestMessage , _ := BuildRequest(pieceBlock)
	size ,id , payload = ParseMsg(requestMessage)
	assert.Equal(t , size ,int32(13), "Request: Size not equal")
	assert.Equal(t , id , int8(6), "Request: Message ID different")
	assert.Equal(t , uint32(payload["index"].(int32)) , pieceBlock.Index , "Request: index field of payload not same") 
	assert.Equal(t , uint32(payload["begin"].(int32)) , pieceBlock.Begin , "Request: begin field of payload not same")
	var length uint32
	binary.Read(payload["length"].(*bytes.Buffer), binary.BigEndian , &length) 
	assert.Equal(t , length , pieceBlock.Length , "Request: length field of payload not same")

	// BuildCancel
	cancelMessage , _ := BuildCancel(pieceBlock)
	size ,id , payload = ParseMsg(cancelMessage)
	assert.Equal(t , size ,int32(13), "Cancel: Size not equal")
	assert.Equal(t , id , int8(8), "Cancel: Message ID different")
	assert.Equal(t , uint32(payload["index"].(int32)) , pieceBlock.Index , "Cancel: index field of payload not same") 
	assert.Equal(t , uint32(payload["begin"].(int32)) , pieceBlock.Begin , "Cancel: begin field of payload not same")
	binary.Read(payload["length"].(*bytes.Buffer), binary.BigEndian , &length) 
	assert.Equal(t , length , pieceBlock.Length , "Cancel: length field of payload not same") 

	// BuildPort
	port1 := uint16(rand.Uint32())
	portMessage , _ := BuildPort(port1)
	size ,id , payload = ParseMsg(portMessage)
	assert.Equal(t , size ,int32(3), "Port: Size not equal")
	assert.Equal(t , id , int8(9), "Port Message ID different")
	var port2 uint16
	binary.Read(payload["payload"].(*bytes.Buffer), binary.BigEndian , &port2)
	assert.Equal(t , p2 , p1 , "Port: length of payload not zero")


}
