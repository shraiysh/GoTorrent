package tracker

// ConnectResponse is struture to hoild details from ConnectResponse
type ConnectResponse struct {
	action        uint32
	transactionId uint32
	connectionId  uint64
}

// AnnounceResponse is structure to hold details from announce request sent to tracker
type AnnounceResponse struct {
	action        uint32
	transactionId uint32
	interval      uint32
	leechers      uint32
	seeders       uint32
	peers         map[uint32]uint16 // peers is a map from IP address to TCP port
}
