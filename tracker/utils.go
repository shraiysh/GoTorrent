package tracker

import (
	"bufio"
	"bytes"
	//"crypto"
	"encoding/binary"
	"fmt"
	"github.com/concurrency-8/parser"
	bencode "github.com/zeebo/bencode"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
	"math/rand"
	"log"
)
 var (
        root string
        torrents []string
        err error
        )
// buildConnReq is the first connection request for tracker
func buildConnReq() []byte {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	binary.Write(writer, binary.BigEndian, uint64(0x41727101980))
	binary.Write(writer, binary.BigEndian, uint32(0))
	binary.Write(writer, binary.BigEndian, getRandomByteArr(4))
	writer.Flush()

	return buffer.Bytes()
}

// respType is for decoding the responses received from socket
func respType(response bytes.Buffer) string {
	action := binary.BigEndian.Uint32(response.Bytes()[0:4])
	if action == 0 {
		return "connect"
	}
	return "announce"
}

// parseConnResp parses the connection request and returns action, transactionID and connectionID
func parseConnResp(response bytes.Buffer) ConnectResponse {
	var connectionResponse ConnectResponse
	responseBytes := response.Bytes()
	connectionResponse.Action = binary.BigEndian.Uint32(responseBytes[0:4])
	connectionResponse.TransactionID = binary.BigEndian.Uint32(responseBytes[4:8])
	connectionResponse.ConnectionID = binary.BigEndian.Uint64(responseBytes[8:])
	return connectionResponse
}

// getrandomByteArr gives a random byte array of specified length
func getRandomByteArr(size uint) []byte {
	temp := make([]byte, size)
	_, err := rand.Read(temp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to generate Crypto random byte array")
	}
	return temp
}

// buildAnnounceReq builds an announce request where we tell the tracker which files we're interested in
func buildAnnounceReq(connectionID uint64, report *ClientStatusReport) (buffer *bytes.Buffer, err error) {
	buffer = new(bytes.Buffer)

	// connection id
	err = binary.Write(buffer, binary.BigEndian, connectionID)
	if err != nil {
		return
	}

	// action
	err = binary.Write(buffer, binary.BigEndian, uint32(1)) // announce req
	if err != nil {
		return
	}

	// transaction id
	err = binary.Write(buffer, binary.BigEndian, getRandomByteArr(4))
	if err != nil {
		return
	}

	// info hash
	var infoHash [20]byte
	copy(infoHash[:], report.TorrentFile.InfoHash)
	err = binary.Write(buffer, binary.BigEndian, infoHash)
	if err != nil {
		return
	}

	// peer id
	err = binary.Write(buffer, binary.BigEndian, []byte(report.PeerID))
	if err != nil {
		return
	}

	// downloaded
	err = binary.Write(buffer, binary.BigEndian, report.Downloaded)
	if err != nil {
		return
	}

	// left
	err = binary.Write(buffer, binary.BigEndian, report.Left)
	if err != nil {
		return
	}

	// uploaded
	err = binary.Write(buffer, binary.BigEndian, report.Uploaded)
	if err != nil {
		return
	}

	// event
	var event uint32
	if report.Event == "" {
		event = 0
	}

	err = binary.Write(buffer, binary.BigEndian, event)
	if err != nil {
		return
	}

	// ip address
	err = binary.Write(buffer, binary.BigEndian, uint32(0))
	if err != nil {
		return
	}

	// key
	err = binary.Write(buffer, binary.BigEndian, uint32(0))
	if err != nil {
		return
	}

	// num want
	err = binary.Write(buffer, binary.BigEndian, int32(-1))
	if err != nil {
		return
	}

	// port
	err = binary.Write(buffer, binary.BigEndian, report.Port)
	if err != nil {
		return
	}

	return
}

// parseAnnounceResp parses necessary details from the announce response sent by tracker
func parseAnnounceResp(response bytes.Buffer) *AnnounceResponse {
	var result AnnounceResponse

	responseBytes := response.Bytes()

	result.Action = binary.BigEndian.Uint32(responseBytes[0:4])
	result.TransactionID = binary.BigEndian.Uint32(responseBytes[4:8])
	result.Interval = binary.BigEndian.Uint32(responseBytes[8:12])
	result.Leechers = binary.BigEndian.Uint32(responseBytes[12:16])
	result.Seeders = binary.BigEndian.Uint32(responseBytes[16:20])

	result.Peers = make([]Peer, (len(responseBytes)-20)/6)

	for i := 20; i+5 < len(responseBytes); i += 6 {
		result.Peers[(i-20)/6].IPAdress = binary.BigEndian.Uint32(responseBytes[i : i+4])
		result.Peers[(i-20)/6].Port = binary.BigEndian.Uint16(responseBytes[i+4 : i+6])
	}

	return &result
}

// getPeersUDP return the list of peers from tracker using UDP urls
func getPeersUDP(u *url.URL, report *ClientStatusReport) (resp *AnnounceResponse, err error) {
	serverAddr, err := net.ResolveUDPAddr("udp", u.Host)
	if err != nil {
		return
	}
	con, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		return
	}
	defer con.Close()

	var connectionID uint64
	for retry := uint(0); retry < uint(8); retry++ {

		err = con.SetDeadline(time.Now().Add(15 * (1 << retry) * time.Second)) // 8 retries
		if err != nil {
			return
		}

		connectionID, err = connectToUDPTracker(con) // get the connection ID
		if err == nil {
			break
		}

		if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
			continue
		}

		if err != nil {
			return
		}

	}

	return getAnnouncementFromUDPTracker(con, connectionID, report)
}

// connnectToUDPTracker send the connection requests and receives connection ID as response
func connectToUDPTracker(con *net.UDPConn) (connectionID uint64, err error) {

	connRequest := buildConnReq()

	_, err = con.Write(connRequest)
	if err != nil {
		return
	}

	respBytes := make([]byte, 16)

	var respLen int
	respLen, err = con.Read(respBytes)
	if err != nil {
		fmt.Println(err)
		return
	}

	if respLen != 16 {
		err = fmt.Errorf("Unexpected response size %d", respLen)
		return
	}

	resp := bytes.NewBuffer(respBytes)
	var resType string
	resType = respType(*resp)

	if resType != "connect" {
		err = fmt.Errorf("Unexpected response action %s", resType)
		return
	}

	connResponse := parseConnResp(*resp)
	connectionID = connResponse.ConnectionID
	return

}

// getAnnouncementFromUDPTracker sends the announce string to UDP tracker and parses the response to get the list of peers
func getAnnouncementFromUDPTracker(con *net.UDPConn, connectionID uint64, report *ClientStatusReport) (resp *AnnounceResponse, err error) {

	announceRequest, err := buildAnnounceReq(connectionID, report)
	if err != nil {
		return
	}

	_, err = con.Write(announceRequest.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}

	respBuffer := new(bytes.Buffer)

	var respLen int
	respBytes := make([]byte, 4096)
	respLen, err = con.Read(respBytes)

	if err != nil {
		return
	}

	if respLen == 0 {
		return
	}

	binary.Write(respBuffer, binary.BigEndian, respBytes[:respLen])

	resp = parseAnnounceResp(*respBuffer)
	return
}

// decodePeerBytes decodes the raw bytes into array of Peers
func (tr *AnnounceResponse) decodePeerBytes() {
	tr.Peers = make([]Peer, len(tr.PeerBytes)/6)

	for i := 0; i+5 < len(tr.PeerBytes); i += 6 {
		tr.Peers[i/6].IPAdress = binary.BigEndian.Uint32(tr.PeerBytes[i : i+4])
		tr.Peers[i/6].Port = binary.BigEndian.Uint16(tr.PeerBytes[i+4 : i+6])
	}
}

// getPeersHTTP returns the list of peers from tracker using HTTP urls
func getPeersHTTP(u *url.URL, report *ClientStatusReport) (tr *AnnounceResponse, err error) {
	uq := u.Query()

	uq.Add("info_hash", report.TorrentFile.InfoHash)
	uq.Add("peer_id", report.PeerID)
	uq.Add("port", strconv.FormatUint(uint64(report.Port), 10))
	uq.Add("uploaded", strconv.FormatUint(report.Uploaded, 10))
	uq.Add("downloaded", strconv.FormatUint(report.Downloaded, 10))
	uq.Add("left", strconv.FormatUint(report.Left, 10))
	uq.Add("compact", "1")

	u.RawQuery = uq.Encode()

	resp, err := http.Get(u.String())

	if err != nil {
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return
	}

	tr = &AnnounceResponse{}
	err = bencode.DecodeBytes(body, tr)
	(*tr).decodePeerBytes()

	return
}

// GetPeers returns the peer list given a valid udp/http announce url
func GetPeers(u *url.URL, report *ClientStatusReport) (tr *AnnounceResponse, err error) {

	switch u.Scheme {
	case "http":
		tr, err = getPeersHTTP(u, report)
	case "udp":
		tr, err = getPeersUDP(u, report)
	default:
		err = fmt.Errorf("Announce url not recognized")
	}

	return
}

// GetClientStatusReport returns the initial report of client
func GetClientStatusReport(torrent parser.TorrentFile, port uint16) (report *ClientStatusReport) {

	report = &ClientStatusReport{}
	report.TorrentFile = torrent
	report.PeerID = string(getRandomByteArr(20))
	report.Left = torrent.Length
	report.Port = port
	report.Event = ""

	return
}
func GetRandomTorrent() (parser.TorrentFile) {
	root = "././test_torrents"
    files, err := ioutil.ReadDir(root)
    if err != nil {
        log.Fatal(err)
    }

    for _, f := range files {
        torrents = append(torrents, f.Name())
    }
    rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
    random_torrent:=root +"/" + torrents[rand.Intn(len(torrents))]
    store, _ := parser.ParseFromFile(random_torrent)
    return store
}
