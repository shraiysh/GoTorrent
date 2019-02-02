package tracker

import "github.com/concurrency-8/parser"

// ConnectResponse is struture to hoild details from ConnectResponse
type ConnectResponse struct {
	Action        uint32
	TransactionID uint32
	ConnectionID  uint64
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
}
