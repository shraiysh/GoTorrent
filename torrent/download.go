package torrent

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"encoding/gob"
	"log"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sethgrid/multibar"

	"github.com/concurrency-8/args"
	"github.com/concurrency-8/parser"
	"github.com/concurrency-8/piece"
	"github.com/concurrency-8/queue"
	"github.com/concurrency-8/tracker"
)

type handler func(tracker.Peer, []byte, net.Conn, *piece.PieceTracker, *queue.Queue, *tracker.ClientStatusReport, Log) error

// MaxTryTracker is the maximum number of times we should try to connect to a tracker
var MaxTryTracker = 3

// MaxTryForTCP is maximum number of times we try to establish a TCP connection
var MaxTryForTCP = 2

// TCPTimeout is the maximum time for which one must wait for connection to a peer
var TCPTimeout time.Duration = 60

// ReadTimeout is the maximum time for which one must wait for the nect message from the peer. If no message arrives till this point, handshake again
var ReadTimeout time.Duration = 60

// MaxTimeoutErrorCount is the maximum number of times we try to read from a conn before restarting it and get a timeout error(ReadTimeout)
var MaxTimeoutErrorCount = 3

var wg sync.WaitGroup

// Log is the logger for current torrent
type Log struct {
	Info  *log.Logger
	Error *log.Logger
}

// DownloadFromFile downloads torrent from path using port
func DownloadFromFile(path string, port int, bar *multibar.ProgressFunc) {

	// Set up logs
	logFolder := filepath.Join("Logs", path)
	os.MkdirAll(logFolder, os.ModePerm)
	logFile, _ := os.Create(filepath.Join(logFolder, "Download.log"))

	Log := Log{}
	Log.Info = log.New(logFile, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
	Log.Error = log.New(logFile, "ERROR ", log.Ldate|log.Ltime|log.Lshortfile)

	torrentFile, err := parser.ParseFromFile(path)
	if err != nil {
		Log.Error.Println("Unable to open torrentfile", err)
		panic(err)
	}
	Log.Info.Println("TorrentFile parsed")

	// Generate client status report
	clientReport := tracker.GetClientStatusReport(torrentFile, uint16(port))

	// Getting peer list from one announce url only for now.
	var announceResp *tracker.AnnounceResponse
	for _, announceURL := range torrentFile.Announce {
		u, err := url.Parse(announceURL)
		if err != nil {
			panic(err)
		}
		Log.Info.Println("Contacting tracker[", announceURL, "] for peer list...")
		count := 0
		for count < MaxTryTracker {
			count++
			announceResp, err = tracker.GetPeers(u, clientReport)
			if err == nil {
				break
			}
			Log.Info.Println("Failed(", err, "). Trying again...")
		}
		if err == nil {
			break
		}
	}

	if announceResp == nil {
		panic("Unable to receive peers! Problem with the torrent or internet")
	}

	pieceTracker := piece.NewPieceTracker(torrentFile)
	if args.ARGS.Resume {
		readGob(torrentFile.Name+"/resume.gob", &pieceTracker.Received, Log)
		readGob(torrentFile.Name+"/resume.gob", &pieceTracker.Requested, Log)
	}

	for i := range pieceTracker.Received {
		temp := true
		for j := range pieceTracker.Received[i] {
			temp = temp && pieceTracker.Received[i][j]
		}

		if !temp {
			for j := range pieceTracker.Received[i] {
				pieceTracker.Received[i][j] = false
			}
		}
	}
	// DownloadFromPeer(announceResp.Peers[0], clientReport, pieceTracker)
	wg.Add(len(announceResp.Peers))
	for _, peer := range announceResp.Peers {
		Log.Info.Println("Spawning peer thread: peer<", peer, ">")
		go func(peer tracker.Peer, clientReport *tracker.ClientStatusReport, pieceTracker *piece.PieceTracker) {
			DownloadFromPeer(peer, clientReport, pieceTracker, Log)
			defer wg.Done()
		}(peer, clientReport, pieceTracker)
	}

	// DownloadFromPeer(announceResp.Peers[0], clientReport, pieceTracker)
	for !pieceTracker.IsDone() {
		over := pieceTracker.PrintPercentageDone()
		(*bar)(over)
		time.Sleep(1 * time.Second)
	}

	wg.Wait()
	pieceTracker.PrintPercentageDone()

	// Close all files
	for _, file := range clientReport.TorrentFile.Files {
		file.FilePointer.Close()
	}

	// Closing log files
	logFile.Close()

	Log.Info.Println("All peer threads finished!")
}

// DownloadFromPeer is a function that handshakes with a peer specified by peer object.
// Concurrently call this function to establish parallel connections to many peers.
func DownloadFromPeer(peer tracker.Peer, report *tracker.ClientStatusReport, pieces *piece.PieceTracker, Log Log) error {
	//safely handle reading using onWholeMessage

	queue := queue.NewQueue(report.TorrentFile)

	exitStatus := 1
	var err error
	for exitStatus == 1 && err == nil && !pieces.IsDone() {
		queue.Choked = true
		conn, err := sendHandshake(peer, report, Log)
		if err != nil {
			break
		}
		exitStatus, err = onWholeMessage(peer, conn, msgHandler, pieces, queue, report, Log)
		if err != nil {
			break
		}
	}

	Log.Info.Println("peer: <", peer, ">: ends!")
	return err
}

func sendHandshake(peer tracker.Peer, report *tracker.ClientStatusReport, Log Log) (conn net.Conn, err error) {
	buffer, err := BuildHandshake(*report)
	if err != nil {
		return nil, err
	}
	peerip := make([]byte, 4)
	binary.BigEndian.PutUint32(peerip, peer.IPAdress)
	service := net.TCPAddr{
		IP:   peerip,
		Port: int(peer.Port),
		Zone: "",
	}
	Log.Info.Println("peer: <", peer, ">: Dialing TCP connection")
	d := net.Dialer{Timeout: TCPTimeout * time.Second}
	err = nil
	count := 0
	for count < MaxTryForTCP {
		count++
		conn, err = d.Dial("tcp", service.String())
		if err != nil {
			Log.Info.Println("peer: <", peer, ">: Unable to set up TCP connection: ", count)
		} else {
			Log.Info.Println("peer: <", peer, ">: Successfully connected to Peer")
			break
		}
	}

	if err != nil {
		Log.Error.Println("peer: <", peer, ">: Could not connect to peer!")
		return nil, err
	}
	Log.Info.Println("peer: <", peer, ">: Handshaking")

	//write the handshake content into the connection.
	_, err = conn.Write(buffer.Bytes())
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func msgHandler(peer tracker.Peer, msg []byte, conn net.Conn, pieces *piece.PieceTracker, queue *queue.Queue, report *tracker.ClientStatusReport, Log Log) error {
	// Log.Info.Println("peer: <", peer, ">: Message:", msg)

	if (len(msg) == int(uint8(msg[0]))+49) && (bytes.Equal(msg[1:20], []byte("BitTorrent protocol"))) {
		Log.Info.Println("peer: <", peer, ">: Handshake successful")
		message, err := BuildInterested()
		if err != nil {
			Log.Info.Println("peer: <", peer, ">: Error", err.Error())
			return err
		}
		conn.Write(message.Bytes())
		// Log.Info.Println("peer: <", peer, ">: Request(", len(message.Bytes()), "): ", message.Bytes())
	} else {

		_, id, payload := ParseMsg(bytes.NewBuffer(msg))

		if id == 0 {
			Log.Info.Println("peer: <", peer, ">: Choke")
			ChokeHandler(peer, conn, pieces, report, Log)
		}
		if id == 1 {
			Log.Info.Println("peer: <", peer, ">: Unchoke")
			UnchokeHandler(peer, conn, pieces, queue, Log)
		}
		if id == 4 {
			Log.Info.Println("peer: <", peer, ">: Have")
			HaveHandler(peer, conn, pieces, queue, payload, Log)
		}
		if id == 5 {
			Log.Info.Println("peer: <", peer, ">: BitField")
			BitFieldHandler(peer, conn, pieces, queue, payload, Log)
		}
		if id == 7 {
			Log.Info.Println("peer: <", peer, ">: Piece")
			PieceHandler(peer, conn, pieces, queue, report, parser.PieceBlock{
				Index: payload["index"].(uint32),
				Begin: payload["begin"].(uint32),
				Bytes: payload["block"].(*bytes.Buffer).Bytes(),
			}, Log)
		}
	}

	return nil

}

// onWholeMessage sends complete messages to callback function
func onWholeMessage(peer tracker.Peer, conn net.Conn, msgHandler handler, pieces *piece.PieceTracker, queue *queue.Queue, report *tracker.ClientStatusReport, Log Log) (status int, err error) {
	buffer := new(bytes.Buffer)
	handshake := true
	resp := make([]byte, 1000)
	msgLen := -1
	count := 0
	for pieces != nil && !pieces.IsDone() {
		conn.SetReadDeadline(time.Now().Add(ReadTimeout * time.Second)) // Setting Read deadline from a connection
		respLen, err := conn.Read(resp)
		//Please look for a better connection handling in the future.
		//Maybe use defer?

		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() && !handshake {
				count++
				if count < MaxTimeoutErrorCount {
					Log.Info.Println("Timeout error - Try again")
					continue
				} else {
					Log.Info.Println("peer: <", peer, ">: Many timeout errors - Peer not responding. Should try to reconnect")
					conn.Close()
					return 1, err
				}
			} else {
				Log.Error.Println("peer: <", peer, ">: Error while reading from connection: ", err)
				Log.Info.Println("peer: <", peer, ">: Restarting connection")
				conn.Close()
				return 1, nil
			}
		}

		binary.Write(buffer, binary.BigEndian, resp[:respLen])

		if handshake {
			Log.Info.Println("peer: <", peer, ">: First message from peer after connection starts - Must be handshake")
			length := uint8((buffer.Bytes())[0])
			msgLen = int(length) + 49
		} else if msgLen == -1 {
			length := binary.BigEndian.Uint32(buffer.Bytes()[0:4])
			// length := uint32((buffer.Bytes())[0:4])
			msgLen = int(length) + 4
			Log.Info.Println("peer: <", peer, ">: New message reception started, len =", msgLen)
			// Log.Info.Println("peer: <", peer, ">: Setting msgLen to", msgLen)
		}

		for len(buffer.Bytes()) >= 4 && msgLen != -1 && len(buffer.Bytes()) >= msgLen {
			Log.Info.Println("peer: <", peer, ">: Message received, msgLen =", msgLen)
			messageBytes := make([]byte, msgLen)
			binary.Read(buffer, binary.BigEndian, messageBytes)
			// Log.Info.Println("peer: <", peer, ">: msgLen:", msgLen)
			msgHandler(peer, messageBytes, conn, pieces, queue, report, Log)
			Log.Info.Println("peer: <", peer, ">: Message handled - setting msgLen = -1")
			msgLen = -1
			handshake = false
			if len(buffer.Bytes()) > 4 {
				length := binary.BigEndian.Uint32(buffer.Bytes()[0:4])
				msgLen = int(length) + 4
				Log.Info.Println("peer: <", peer, ">: New message was in previous one - msgLen =", msgLen)
				// Log.Info.Println("peer: <", peer, ">: Setting msgLen to", msgLen)
			}
		}
	}
	return 0, nil
}

// ChokeHandler handles choking protocol
func ChokeHandler(peer tracker.Peer, conn net.Conn, pieces *piece.PieceTracker, report *tracker.ClientStatusReport, Log Log) {
	Log.Info.Println("peer:<", peer, ">: Choke")
	if pieces != nil && pieces.IsDone() {
		Log.Info.Println("All pieces done. Closing connection.")
		conn.Close()
	} else if report != nil {
		Log.Info.Println("peer: <", peer, ">: Handshaking again")
		// time.Sleep(2 * time.Second) // Sleep for 2 seconds and try handshaking again
		handshake, err := BuildHandshake(*report)
		if err != nil {
			panic("Problem with the torrentFile")
		} else {
			conn.Write(handshake.Bytes())
		}
	}
}

// UnchokeHandler handles unchoking protocol
func UnchokeHandler(peer tracker.Peer, conn net.Conn, pieces *piece.PieceTracker, queue *queue.Queue, Log Log) {
	if queue.Choked && queue.Length() != 0 {
		Log.Info.Println("peer:<", peer, "> Unchoke: queue was choked, but queue was non-empty")
		queue.Choked = false
		Log.Info.Println("peer:<", peer, ">: Requesting next piece")
		RequestPiece(peer, conn, pieces, queue, Log)
	} else if queue.Choked {
		Log.Info.Println("peer:<", peer, ">: Unchoke - Queue empty and choked - Sending interested")
		queue.Choked = false
		message, _ := BuildInterested()
		// if err != nil {
		// 	Log.Info.Println("peer: <", peer, ">: Error", err.Error())
		// 	return err
		// }
		if conn != nil {
			conn.Write(message.Bytes())
		}
	}
	// Log.Info.Println("peer: <", peer, ">: RequestPiece : Called from Unchokehandler")
}

// HaveHandler handles Have protocol
func HaveHandler(peer tracker.Peer, conn net.Conn, pieces *piece.PieceTracker, queue *queue.Queue, payload Payload, Log Log) (pieceIndex uint32, err error) {
	binary.Read(payload["payload"].(*bytes.Buffer), binary.BigEndian, &pieceIndex)
	queueempty := (queue.Length() == 0)
	err = queue.Enqueue(pieceIndex)
	if err != nil {
		return
	}
	if queueempty {
		Log.Info.Println("peer: <", peer, ">: HaveHandler: Queue was empty. Requesting pieces.")
		err = RequestPiece(peer, conn, pieces, queue, Log)
	}
	return
}

// BitFieldHandler handles bitfield protocol
func BitFieldHandler(peer tracker.Peer, conn net.Conn, pieces *piece.PieceTracker, queue *queue.Queue, payload Payload, Log Log) (err error) {
	queueempty := (queue.Length() == 0)
	msg := payload["payload"]
	for i, bytevalue := range msg.(*bytes.Buffer).Bytes() {
		for j := 7; j >= 0; j-- {
			if 1 == bytevalue&1 {
				err = queue.Enqueue(uint32(i*8 + j))
			}
			bytevalue = bytevalue >> 1
		}
	}
	if queueempty {
		Log.Info.Println("peer: <", peer, ">: BitFieldHandler: Queue was empty. Requesting pieces")
		err = RequestPiece(peer, conn, pieces, queue, Log)
	}

	return
}

// PieceHandler - TODO Write comment
func PieceHandler(peer tracker.Peer, conn net.Conn, pieces *piece.PieceTracker, queue *queue.Queue, report *tracker.ClientStatusReport, pieceResp parser.PieceBlock, Log Log) {
	pieces.AddReceived(pieceResp)

	Log.Info.Println("peer: <", peer, ">: Received piece[", pieceResp.Index, "] [", pieceResp.Begin/parser.BLOCK_LEN, "]")
	report.Data[pieceResp.Index].Blocks[pieceResp.Begin/parser.BLOCK_LEN] = pieceResp

	toSHA1 := func(data []byte) []byte {
		hash := sha1.New()
		hash.Write(data)
		return hash.Sum(nil)
	}
	var piece []byte
	if pieces.PieceIsDone(pieceResp.Index) {
		for _, i := range report.Data[pieceResp.Index].Blocks {
			if len(i.Bytes) != 0 {
				piece = append(piece, i.Bytes...)
			}
		}
		// numZeros := int(report.TorrentFile.PieceLength) - len(piece)
		// for i := 0; i < numZeros; i++ {
		// 	piece = append(piece, byte(0))
		// }
		same := true
		expected := report.TorrentFile.Piece[pieceResp.Index*20 : (pieceResp.Index+1)*20]
		actual := toSHA1(piece)
		for i := range expected {
			same = same && expected[i] == actual[i]
		}
		if !same {
			Log.Error.Println("peer: <", peer, ">: SHA do not match for piece:", pieceResp.Index)
			Log.Error.Println("peer: <", peer, ">: Expected:\t", report.TorrentFile.Piece[pieceResp.Index*20:(pieceResp.Index+1)*20])
			Log.Error.Println("peer: <", peer, ">: Actual:\t", toSHA1(piece))
			report.Data[pieceResp.Index].Blocks[pieceResp.Begin/parser.BLOCK_LEN] = parser.PieceBlock{}

			pieces.Reset(pieceResp.Index)
			queue.Enqueue(pieceResp.Index)
			Log.Info.Println("peer: <", peer, ">: Reset queue and pieceTracker for", pieceResp.Index)
			RequestPiece(peer, conn, pieces, queue, Log)
			return
		}
		Log.Info.Println("peer: <", peer, ">: Piece[", pieceResp.Index, "] downloaded SUCCESSFULLY!")
	}

	// Log.Info.Println("Bytes Received : ", pieceResp.Bytes)

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
	Log.Info.Println("peer: <", peer, ">: Writing block to file ", file.Name())
	file.WriteAt(pieceResp.Bytes, int64(offsetInFile))
	if args.ARGS.ResumeCapability {
		writeGob(report.TorrentFile.Name+"/resume.gob", pieces.Received, Log)
	}

	// file.Sync()
	pieces.PrintPercentageDone()

	if pieces.IsDone() {
		for _, file := range report.TorrentFile.Files {
			defer file.FilePointer.Close()
		}
		Log.Info.Println("peer: <", peer, ">: Done")
		conn.Close()
	} else {
		Log.Info.Println("peer<", peer, " >: Called from piecehandler")
		RequestPiece(peer, conn, pieces, queue, Log)
	}
}

// RequestPiece requests a piece
func RequestPiece(peer tracker.Peer, conn net.Conn, pieces *piece.PieceTracker, queue *queue.Queue, Log Log) (err error) {
	if queue.Choked {
		Log.Error.Println("peer: <", peer, ">: Queue is choked")
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

		pieces.Lock.Lock()
		if pieces.Needed(pieceBlock) {
			pieces.AddRequested(pieceBlock)
			pieces.Lock.Unlock()
			Log.Info.Println("peer: <", peer, ">: Requesting piece[", pieceBlock.Index, "][", pieceBlock.Begin/parser.BLOCK_LEN, "]")
			message, err := BuildRequest(pieceBlock)

			if err != nil {
				queue.Enqueue(pieceBlock.Index)
				break
			}
			_, err = conn.Write(message.Bytes())

			if err != nil {
				Log.Info.Println("peer: <", peer, ">:", err.Error())
				queue.Enqueue(pieceBlock.Index)
				break
			}
			break
		} else {
			pieces.Lock.Unlock()
		}
	}
	return
}

func writeGob(filePath string, object interface{}, Log Log) error {
	file, err := os.Create(filePath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

func readGob(filePath string, object interface{}, Log Log) error {
	file, err := os.Open(filePath)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}
