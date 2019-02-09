package torrent

import(
	"bytes"
	"encoding/binary"
	"github.com/concurrency-8/tracker"
)

/** 
  * BuildHandshake returns a pointer to a buffer.
  * Buffer looks like: 
      uint8		: pstrlen	- Length of pstr
      [pstrlen]byte	: pstr		- pstr, the string identifier of the protocol
      [8]byte		: reserved	- 8 reserved bytes
      [20]byte		: infohash	- SHA1 hash of the info key in the metainfo file. Same as the info hash transmitted in tracker requests
      [20]byte		: peerID	- 20 byte unique ID for the client. Usually the same peerID transmitted in tracker requests
  * In version 1.0 of the BitTorrent protocol, pstrlen = 19, and pstr = "BitTorrent protocol"
*/
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

/**
  * BuildKeepAlive returns pointer to an empty buffer (4 bytes)
  */

func BuildKeepAlive() (keepAlive *bytes.Buffer) {
	keepAlive = bytes.NewBuffer(make([]byte, 4))

	return
}

/**
  * BuildChoke returns pointer to a buffer.
  * Buffer looks like:
      uint32	: length	- Length of remaining part(message) = 1
      uint8	: message	- For choke, message = 0
*/

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

/**
  * BuildUnchoke returns pointer to a buffer.
  * Buffer looks like:
      uint32	: length	- Length of remaining part(message) = 1
      uint8	: message	- For unchoke, message = 1
*/
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

/**
  * BuildInterested returns pointer to a buffer.
  * Buffer looks like:
      uint32	: length	- Length of remaining part(message) = 1
      uint8	: message	- For interested, message = 2
*/
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

/**
  * BuildUninterested returns pointer to a buffer.
  * Buffer looks like:
      uint32	: length	- Length of remaining part(message) = 1
      uint8	: message	- For uninterested, message = 3
*/
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

/**
  * BuildHave returns pointer to a buffer. This takes uint32 payload as an argument
  * Buffer looks like:
      uint32	: length	- Length of remaining part(message) = 5
      uint8	: messageType	- for payload message = 4
      uint32	: peice index	- payload
*/
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
