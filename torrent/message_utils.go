package torrent

import (
	"bytes"
	"encoding/binary"
	"github.com/concurrency-8/queue"
	"github.com/concurrency-8/tracker"
)

// BuildHandshake returns a pointer to a buffer.
// Buffer looks like:
//	uint8		: pstrlen	- Length of pstr
//	[pstrlen]byte	: pstr		- pstr, the string identifier of the protocol
//	[8]byte		: reserved	- 8 reserved bytes
//	[20]byte	: infohash	- SHA1 hash of the info key in the metainfo file. Same as the info hash transmitted in tracker requests
//	[20]byte	: peerID	- 20 byte unique ID for the client. Usually the same peerID transmitted in tracker requests
// In version 1.0 of the BitTorrent protocol, pstrlen = 19, and pstr = "BitTorrent protocol"
func BuildHandshake(report tracker.ClientStatusReport) (handshake *bytes.Buffer, err error) {
	handshake = new(bytes.Buffer)

	// pstrlen
	if err = binary.Write(handshake, binary.BigEndian, uint8(19)); err != nil {
		return
	}

	// pstr
	if err = binary.Write(handshake, binary.BigEndian, []byte("BitTorrent protocol")); err != nil {
		return
	}

	// reserved
	if err = binary.Write(handshake, binary.BigEndian, uint64(0)); err != nil {
		return
	}

	// infohash
	var infohashFromFile [20]byte
	copy(infohashFromFile[:], report.TorrentFile.InfoHash)
	if err = binary.Write(handshake, binary.BigEndian, infohashFromFile); err != nil {
		return
	}

	// peerID
	if err = binary.Write(handshake, binary.BigEndian, []byte(report.PeerID)); err != nil {
		return
	}

	return
}

// BuildKeepAlive returns pointer to an empty buffer (4 bytes)
func BuildKeepAlive() (keepAlive *bytes.Buffer) {
	keepAlive = bytes.NewBuffer(make([]byte, 4))

	return
}

// BuildChoke returns pointer to a buffer.
// Buffer looks like:
//	uint32	: length	- Length of remaining part(message) = 1
//	uint8	: messageType	- For choke, messageType = 0
func BuildChoke() (choke *bytes.Buffer, err error) {
	choke = new(bytes.Buffer)

	if err = binary.Write(choke, binary.BigEndian, uint32(1)); err != nil {
		return
	}

	if err = binary.Write(choke, binary.BigEndian, uint8(0)); err != nil {
		return
	}

	return
}

// BuildUnchoke returns pointer to a buffer.
// Buffer looks like:
//	uint32	: length	- Length of remaining part(message) = 1
//	uint8	: messageType	- For unchoke, messageType = 1
func BuildUnchoke() (unchoke *bytes.Buffer, err error) {
	unchoke = new(bytes.Buffer)

	if err = binary.Write(unchoke, binary.BigEndian, uint32(1)); err != nil {
		return
	}

	if err = binary.Write(unchoke, binary.BigEndian, uint8(1)); err != nil {
		return
	}

	return
}

// BuildInterested returns pointer to a buffer.
// Buffer looks like:
//	uint32	: length	- Length of remaining part(message) = 1
//	uint8	: messageType	- For interested, messageType = 2
func BuildInterested() (interested *bytes.Buffer, err error) {
	interested = new(bytes.Buffer)

	if err = binary.Write(interested, binary.BigEndian, uint32(1)); err != nil {
		return
	}

	if err = binary.Write(interested, binary.BigEndian, uint8(2)); err != nil {
		return
	}

	return
}

// BuildUninterested returns pointer to a buffer.
// Buffer looks like:
//	uint32	: length	- Length of remaining part(message) = 1
//	uint8	: messageType	- For uninterested, messageType = 3
func BuildUninterested() (uninterested *bytes.Buffer, err error) {
	uninterested = new(bytes.Buffer)

	if err = binary.Write(uninterested, binary.BigEndian, uint32(1)); err != nil {
		return
	}

	if err = binary.Write(uninterested, binary.BigEndian, uint8(3)); err != nil {
		return
	}

	return
}

// BuildHave returns pointer to a buffer. This takes uint32 payload(piece index) as an argument
// Buffer looks like:
//	uint32	: length	- Length of remaining part(message) = 5
//	uint8	: messageType	- for have, messageType = 4
//	uint32	: piece index	- payload
func BuildHave(payload uint32) (have *bytes.Buffer, err error) {
	have = new(bytes.Buffer)

	if err = binary.Write(have, binary.BigEndian, uint32(5)); err != nil {
		return
	}

	if err = binary.Write(have, binary.BigEndian, uint8(4)); err != nil {
		return
	}

	if err = binary.Write(have, binary.BigEndian, uint32(payload)); err != nil {
		return
	}

	return
}

// BuildRequest returns pointer to a buffer. This takes queue.PieceBlock as an argument
//	uint32	: length	- Length of remaining part(message) = 13
//	uint8	: messageType	- for request, message = 5
//	uint32	: piece index	- queue.PieceBlock.Index for payload
//	uint32	: piece begin	- queue.PieceBlock.Begin for payload
//	uint32	: piece length	- queue.PieceBlock.Length for payload
func BuildRequest(payload queue.PieceBlock) (request *bytes.Buffer, err error) {
	request = new(bytes.Buffer)

	// Length of message
	if err = binary.Write(request, binary.BigEndian, uint32(13)); err != nil {
		return
	}
	// message type
	if err = binary.Write(request, binary.BigEndian, uint8(6)); err != nil {
		return
	}
	// piece index
	if err = binary.Write(request, binary.BigEndian, uint32(payload.Index)); err != nil {
		return
	}
	// piece begin
	if err = binary.Write(request, binary.BigEndian, uint32(payload.Begin)); err != nil {
		return
	}
	// piece length
	if err = binary.Write(request, binary.BigEndian, uint32(payload.Length)); err != nil {
		return
	}

	return
}

// BuildPiece returns pointer to a buffer having the piece. Takes the queue.PieceBlock object as an arg
//	uint32	: length	- length of remaining part (message) = payload length + 9
//	uint8	: messageType	- for piece, type = 7
//	uint32	: piece index	- queue.PieceBlock.Index for payload
//	uint32	: piece begin	- queue.PieceBlock.Begin for payload
//	[]byte	: piece		- the data of the piece, queue.PieceBlock.Block for payload
// func BuildPiece(payload queue.PieceBlock) (piece *bytes.Buffer, err error) {
// 	piece = new(bytes.Buffer)

// 	// Length of message (Has the piece)
// 	if err = binary.Write(piece, binary.BigEndian, uint32(len(payload.Block.Bytes())+9)); err != nil {
// 		return
// 	}

// 	// Message type
// 	if err = binary.Write(piece, binary.BigEndian, uint8(7)); err != nil {
// 		return
// 	}

// 	// piece index
// 	if err = binary.Write(piece, binary.BigEndian, uint32(payload.Index)); err != nil {
// 		return
// 	}

// 	// piece begin
// 	if err = binary.Write(piece, binary.BigEndian, uint32(payload.Begin)); err != nil {
// 		return
// 	}

// 	// piece
// 	if err = binary.Write(piece, binary.BigEndian, payload.Block.Bytes()); err != nil {
// 		return
// 	}

// 	return
// }

// BuildCancel returns pointer to a buffer. Takes queue.PieceBlock object as arg
//	uint32	: length	- Length of the remaining message = 13
//	uint8	: messageType	- for cancel, messageType = 8
//	uint32	: piece index	- queue.PieceBlock.Index for payload
//	uint32	: piece begin	- queue.PieceBlock.Begin for payload
//	uint32	: piece length	- queue.PieceBlock.Length for payload
func BuildCancel(payload queue.PieceBlock) (cancelBuf *bytes.Buffer, err error) {
	cancelBuf = new(bytes.Buffer)

	// Length of Message
	if err = binary.Write(cancelBuf, binary.BigEndian, uint32(13)); err != nil {
		return
	}

	// Message type - Cancel
	if err = binary.Write(cancelBuf, binary.BigEndian, uint8(8)); err != nil {
		return
	}

	// piece index
	if err = binary.Write(cancelBuf, binary.BigEndian, payload.Index); err != nil {
		return
	}

	// piece begin
	if err = binary.Write(cancelBuf, binary.BigEndian, payload.Begin); err != nil {
		return
	}

	// piece length
	if err = binary.Write(cancelBuf, binary.BigEndian, payload.Length); err != nil {
		return
	}

	return
}

// BuildPort returns a pointer to a buffer. Takes uint16 port as arg
//	uint32	: length	- length of remaining message = 3
//	uint8	: messageType	- for port, messageType = 9
//	uint16	: port		- the argument, port
func BuildPort(port uint16) (portBuf *bytes.Buffer, err error) {
	portBuf = new(bytes.Buffer)

	// Length of message
	if err = binary.Write(portBuf, binary.BigEndian, uint32(3)); err != nil {
		return
	}
	// Message type = 9
	if err = binary.Write(portBuf, binary.BigEndian, uint8(9)); err != nil {
		return
	}
	// listen-port
	if err = binary.Write(portBuf, binary.BigEndian, uint16(port)); err != nil {
		return
	}
	return
}

// ParseMsg parses a message
func ParseMsg(msg *bytes.Buffer) (size int32 , id int8 , payload Payload){

	binary.Read(msg,binary.BigEndian,&size);

	if size > 0 {
		binary.Read(msg , binary.BigEndian , &id);
	}

	if ( id == 6 || id ==7 || id==8){
		rest = bytes.NewBuffer(msg.Bytes()[8:])
		binary.Read(msg,binary.BigEndian,&payload.Index)
		binary.Read(msg,binary.BigEndian,&payload.Begin)

		if id == 7 {
			payload.Block = rest
		}else{
			payload.Length = rest 
		}
	}

	return

}
