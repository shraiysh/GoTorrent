package torrent

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"github.com/concurrency-8/parser"
	"github.com/concurrency-8/piece"
	"github.com/concurrency-8/queue"
	"github.com/concurrency-8/tracker"
	"net"
	"net/url"
	"sync"
)

type handler func(tracker.Peer, []byte, net.Conn, *piece.PieceTracker, *queue.Queue, *tracker.ClientStatusReport) error

// MAX_TRY is the maximum number of times we should try to connect to a tracker
var MAX_TRY int = 5

var wg sync.WaitGroup

// DownloadFromFile downloads torrent from path using port
func DownloadFromFile(path string, port int) {
	torrentFile, err := parser.ParseFromFile(path)
	if err != nil {
		panic(err)
	}
	fmt.Println("TorrentFile parsed")

	// Generate client status report
	clientReport := tracker.GetClientStatusReport(torrentFile, uint16(port))

	// Getting peer list from one announce url only for now.
	var announceResp *tracker.AnnounceResponse
	for _, announceUrl := range torrentFile.Announce {
		u, err := url.Parse(announceUrl)
		if err != nil {
			panic(err)
		}
		fmt.Println("Contacting tracker[", announceUrl, "] for peer list...")
		count := 0
		for count < MAX_TRY {
			count++
			announceResp, err = tracker.GetPeers(u, clientReport)
			if err == nil {
				break
			}
			fmt.Println("Failed(", err, "). Trying again...")
		}
		if err == nil {
			break
		}
	}

	if announceResp == nil {
		panic("Unable to receive peers! Problem with the torrent or internet")
	}

	pieceTracker := piece.NewPieceTracker(torrentFile)
	wg.Add(len(announceResp.Peers))
	for _, peer := range announceResp.Peers {
		fmt.Println("Spawning peer thread: peer<", peer, ">")
		go DownloadFromPeer(peer, clientReport, pieceTracker)
	}

	// DownloadFromPeer(announceResp.Peers[0], clientReport, pieceTracker)

	wg.Wait()
	fmt.Println("All peer threads finished!")
}

// DownloadFromPeer is a function that handshakes with a peer specified by peer object.
// Concurrently call this function to establish parallel connections to many peers.
func DownloadFromPeer(peer tracker.Peer, report *tracker.ClientStatusReport, pieces *piece.PieceTracker) error {
	fmt.Println("peer: <", peer, ">: Handshaking")
	buffer, err := BuildHandshake(*report)
	if err != nil {
		return err
	}
	peerip := make([]byte, 4)
	binary.BigEndian.PutUint32(peerip, peer.IPAdress)
	service := net.TCPAddr{
		IP:   peerip,
		Port: int(peer.Port),
		Zone: "",
	}
	conn, err := net.Dial("tcp", service.String())
	if err != nil {
		return err
	}
	//write the handshake content into the connection.
	_, err = conn.Write(buffer.Bytes())
	if err != nil {
		return err
	}
	//safely handle reading using onWholeMessage

	queue := queue.NewQueue(report.TorrentFile)
	// for !pieces.IsDone() {
	onWholeMessage(peer, conn, msgHandler, pieces, queue, report)
	// }

	fmt.Println("peer: <", peer, ">: ends!")
	return err
}
func msgHandler(peer tracker.Peer, msg []byte, conn net.Conn, pieces *piece.PieceTracker, queue *queue.Queue, report *tracker.ClientStatusReport) error {
	// fmt.Println("peer: <", peer, ">: Message:", msg)

	if (len(msg) == int(uint8(msg[0]))+49) && (bytes.Equal(msg[1:20], []byte("BitTorrent protocol"))) {
		fmt.Println("peer: <", peer, ">: Handshake successful")
		message, err := BuildInterested()
		if err != nil {
			// fmt.Println("peer: <", peer, ">: Error", err.Error())
			return err
		}
		conn.Write(message.Bytes())
		// fmt.Println("peer: <", peer, ">: Request(", len(message.Bytes()), "): ", message.Bytes())
	} else {

		_, id, payload := ParseMsg(bytes.NewBuffer(msg))

		if id == 0 {
			// fmt.Println("peer: <", peer, ">: Choke")
			ChokeHandler(peer, conn)
		}
		if id == 1 {
			// fmt.Println("peer: <", peer, ">: Unchoke")
			UnchokeHandler(peer, conn, pieces, queue)
		}
		if id == 4 {
			// fmt.Println("peer: <", peer, ">: Have")
			HaveHandler(peer, conn, pieces, queue, payload)
		}
		if id == 5 {
			// fmt.Println("peer: <", peer, ">: BitField")
			BitFieldHandler(peer, conn, pieces, queue, payload)
		}
		if id == 7 {
			// fmt.Println("peer: <", peer, ">: Piece")
			PieceHandler(peer, conn, pieces, queue, report, parser.PieceBlock{
				Index: payload["index"].(uint32),
				Begin: payload["begin"].(uint32),
				Bytes: payload["block"].(*bytes.Buffer).Bytes(),
			})
		}
	}

	return nil

}

// onWholeMessage sends complete messages to callback function
func onWholeMessage(peer tracker.Peer, conn net.Conn, msgHandler handler, pieces *piece.PieceTracker, queue *queue.Queue, report *tracker.ClientStatusReport) error {
	buffer := new(bytes.Buffer)
	handshake := true
	resp := make([]byte, 1000)
	for {
		respLen, err := conn.Read(resp)
		//Please look for a better connection handling in the future.
		//Maybe use defer?

		if err != nil {
			conn.Close()
			return err
		}

		binary.Write(buffer, binary.BigEndian, resp[:respLen])

		var msgLen int
		if handshake {
			length := uint8((buffer.Bytes())[0])
			msgLen = int(length) + 49
		} else {
			length := binary.BigEndian.Uint32(buffer.Bytes()[0:4])
			// length := uint32((buffer.Bytes())[0:4])
			msgLen = int(length) + 4
			// fmt.Println("peer: <", peer, ">: Setting msgLen to", msgLen)
		}

		for len(buffer.Bytes()) >= 4 && len(buffer.Bytes()) >= msgLen {
			messageBytes := make([]byte, msgLen)
			binary.Read(buffer, binary.BigEndian, messageBytes)
			// fmt.Println("peer: <", peer, ">: msgLen:", msgLen)
			msgHandler(peer, messageBytes, conn, pieces, queue, report)
			handshake = false
			if len(buffer.Bytes()) > 4 {
				length := binary.BigEndian.Uint32(buffer.Bytes()[0:4])
				msgLen = int(length) + 4
				// fmt.Println("peer: <", peer, ">: Setting msgLen to", msgLen)
			}

		}
	}
}

// ChokeHandler handles choking protocol
func ChokeHandler(peer tracker.Peer, conn net.Conn) {
	fmt.Println("peer:<", peer, ">: Choke: close")
	conn.Close()
}

// UnchokeHandler handles unchoking protocol
func UnchokeHandler(peer tracker.Peer, conn net.Conn, pieces *piece.PieceTracker, queue *queue.Queue) {
	queue.Choked = false
	RequestPiece(peer, conn, pieces, queue)
}

// HaveHandler handles Have protocol
func HaveHandler(peer tracker.Peer, conn net.Conn, pieces *piece.PieceTracker, queue *queue.Queue, payload Payload) {
	return
}

// BitFieldHandler handles bitfield protocol
func BitFieldHandler(peer tracker.Peer, conn net.Conn, pieces *piece.PieceTracker, queue *queue.Queue, payload Payload) {
	// RequestPiece(conn, pieces, queue)
	for i := range pieces.Requested {
		queue.Enqueue(uint32(i))
	}
	return
}

// PieceHandler - TODO Write comment
func PieceHandler(peer tracker.Peer, conn net.Conn, pieces *piece.PieceTracker, queue *queue.Queue, report *tracker.ClientStatusReport, pieceResp parser.PieceBlock) {
	pieces.AddReceived(pieceResp)

	fmt.Println("peer: <", peer, ">: Received piece[", pieceResp.Index, "] [", pieceResp.Begin/parser.BLOCK_LEN, "]")

	offsetInFile := uint64(pieceResp.Index)*uint64(report.TorrentFile.PieceLength) + uint64(pieceResp.Begin)
	file := report.TorrentFile.Files[0].FilePointer
	for key, value := range report.TorrentFile.Files {
		if offsetInFile > value.Length {
			offsetInFile -= value.Length
			file = report.TorrentFile.Files[key+1].FilePointer
		} else {
			break
		}
	}
	report.Data[pieceResp.Index].Blocks[pieceResp.Begin/parser.BLOCK_LEN] = pieceResp
	fmt.Println("peer: <", peer, ">: Writing block to file ", file.Name())
	file.WriteAt(pieceResp.Bytes, int64(offsetInFile))
	file.Sync()

	toSHA1 := func(data []byte) []byte {
		hash := sha1.New()
		hash.Write(data)
		return hash.Sum(nil)
	}
	var piece []byte
	if pieces.PieceIsDone(pieceResp.Index) {
		for _, i := range report.Data[pieceResp.Index].Blocks {
			piece = append(piece, i.Bytes...)
		}

		same := true
		expected := report.TorrentFile.Piece[pieceResp.Index*20 : (pieceResp.Index+1)*20]
		actual := toSHA1(piece)
		for i := range expected {
			same = same && expected[i] == actual[i]
		}
		if !same {
			fmt.Println("peer: <", peer, ">: Expected:\t", report.TorrentFile.Piece[pieceResp.Index*20:(pieceResp.Index+1)*20])
			fmt.Println("peer: <", peer, ">: Actual:\t", toSHA1(piece))
			panic("Error downloading! SHA don't match")

		} else {
			fmt.Println("peer: <", peer, ">: Piece[", pieceResp.Index, "] downloaded SUCCESSFULLY!")
		}
	}

	if pieces.IsDone() {
		for _, file := range report.TorrentFile.Files {
			defer file.FilePointer.Close()
		}
		fmt.Println("peer: <", peer, ">: Done")
		conn.Close()
	} else {
		RequestPiece(peer, conn, pieces, queue)
	}
}

// RequestPiece requests a piece
func RequestPiece(peer tracker.Peer, conn net.Conn, pieces *piece.PieceTracker, queue *queue.Queue) (err error) {
	if queue.Choked {
		err = fmt.Errorf("peer: <", peer, ">: Queue is choked")
		return
	}

	for queue.Length() > 0 {
		pieceBlock, err := queue.Peek()

		if err != nil {
			break
		}

		err = queue.Dequeue()

		if err != nil {
			break
		}

		if pieces.Needed(pieceBlock) {
			fmt.Println("peer: <", peer, ">: Requesting piece[", pieceBlock.Index, "][", pieceBlock.Begin/parser.BLOCK_LEN, "]")
			message, err := BuildRequest(pieceBlock)

			if err != nil {
				break
			}
			_, err = conn.Write(message.Bytes())

			if err != nil {
				break
			}

			pieces.AddRequested(pieceBlock)
			break
		}
	}
	return
}
