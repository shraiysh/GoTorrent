package torrent

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/concurrency-8/parser"
	"github.com/concurrency-8/piece"
	"github.com/concurrency-8/queue"
	"github.com/concurrency-8/tracker"
	"net"
	"os"
)

type handler func([]byte, net.Conn, *piece.PieceTracker, *queue.Queue, *tracker.ClientStatusReport) error

// Download is a function that handshakes with a peer specified by peer object.
// Concurrently call this function to establish parallel connections to many peers.
func Download(peer tracker.Peer, report *tracker.ClientStatusReport, pieces *piece.PieceTracker) error {
	fmt.Println("torrent::Download")
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
	fmt.Println("Request(", len(buffer.Bytes()), "): ", buffer.Bytes())
	if err != nil {
		return err
	}
	//safely handle reading using onWholeMessage

	queue := queue.NewQueue(report.TorrentFile)
	// for !pieces.IsDone() {
	onWholeMessage(conn, msgHandler, pieces, queue, report)
	// }
	return err
}
func msgHandler(msg []byte, conn net.Conn, pieces *piece.PieceTracker, queue *queue.Queue, report *tracker.ClientStatusReport) error {
	fmt.Println("Message:", msg)

	if (len(msg) == int(uint8(msg[0]))+49) && (bytes.Equal(msg[1:20], []byte("BitTorrent protocol"))) {
		fmt.Println("Handshake")
		message, err := BuildInterested()
		if err != nil {
			fmt.Println("Error", err.Error())
			return err
		}
		conn.Write(message.Bytes())
		fmt.Println("Request(", len(message.Bytes()), "): ", message.Bytes())
	} else {

		_, id, payload := ParseMsg(bytes.NewBuffer(msg))

		if id == 0 {
			fmt.Println("Choke")
			ChokeHandler(conn)
		}
		if id == 1 {
			fmt.Println("Unchoke")
			UnchokeHandler(conn, pieces, queue)
		}
		if id == 4 {
			fmt.Println("Have")
			HaveHandler(conn, pieces, queue, payload)
		}
		if id == 5 {
			fmt.Println("BitField")
			BitFieldHandler(conn, pieces, queue, payload)
		}
		if id == 7 {
			fmt.Println("Piece")
			fmt.Println(payload)
			// PieceHandler(conn, pieces, queue, report.TorrentFile, payload)
		}
	}

	return nil

}

// onWholeMessage sends complete messages to callback function
func onWholeMessage(conn net.Conn, msgHandler handler, pieces *piece.PieceTracker, queue *queue.Queue, report *tracker.ClientStatusReport) error {
	buffer := new(bytes.Buffer)
	handshake := true
	resp := make([]byte, 100)
	for {
		respLen, err := conn.Read(resp)
		//Please look for a better connection handling in the future.
		//Maybe use defer?
		fmt.Println("Response(", respLen, "): ", resp[:respLen])
		bufio.NewReader(os.Stdin).ReadBytes('\n')

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
			fmt.Println("Setting msgLen to", msgLen)
		}

		for len(buffer.Bytes()) >= 4 && len(buffer.Bytes()) >= msgLen {
			messageBytes := make([]byte, msgLen)
			binary.Read(buffer, binary.BigEndian, messageBytes)
			// fmt.Println("msgLen:", msgLen)
			// fmt.Println("message:", messageBytes)
			msgHandler(messageBytes, conn, pieces, queue, report)
			handshake = false
			if len(buffer.Bytes()) > 0 {
				length := int32((buffer.Bytes())[0])
				msgLen = int(length) + 4
			}

		}
	}
}

// ChokeHandler handles choking protocol
func ChokeHandler(conn net.Conn) {
	conn.Close()
}

// UnchokeHandler handles unchoking protocol
func UnchokeHandler(conn net.Conn, pieces *piece.PieceTracker, queue *queue.Queue) {
	queue.Choked = false
	RequestPiece(conn, pieces, queue)
}

// HaveHandler handles Have protocol
func HaveHandler(conn net.Conn, pieces *piece.PieceTracker, queue *queue.Queue, payload Payload) {
	return
}

// BitFieldHandler handles bitfield protocol
func BitFieldHandler(conn net.Conn, pieces *piece.PieceTracker, queue *queue.Queue, payload Payload) {
	// RequestPiece(conn, pieces, queue)
	return
}

// PieceHandler - TODO Write comment
func PieceHandler(conn net.Conn, pieces *piece.PieceTracker, queue *queue.Queue, torrent *parser.TorrentFile, pieceResp *parser.PieceBlock) {
	pieces.AddReceived(*pieceResp)
	fmt.Println(pieceResp.Bytes)

	if pieces.IsDone() {
		fmt.Println("Done")
		conn.Close()
	} else {
		RequestPiece(conn, pieces, queue)
	}
}

// RequestPiece requests a piece
func RequestPiece(conn net.Conn, pieces *piece.PieceTracker, queue *queue.Queue) (err error) {
	fmt.Println("RequestPiece")
	if (*queue).Choked {
		fmt.Println("Queue is Choked!")
		err = fmt.Errorf("Queue is choked")
		return
	}
	fmt.Println("Outside if")

	for queue.Length() > 0 {
		fmt.Println("Queue not empty")
		pieceBlock, err := queue.Peek()

		if err != nil {
			break
		}

		err = queue.Dequeue()

		if err != nil {
			break
		}

		fmt.Println("Checking Piece:", pieceBlock.Index)

		if pieces.Needed(pieceBlock) {
			fmt.Println("Requesting piece:", pieceBlock.Index)
			message, err := BuildRequest(pieceBlock)

			if err != nil {
				break
			}
			_, err = conn.Write(message.Bytes())
			fmt.Println("Request(", len(message.Bytes()), "): ", message.Bytes())

			if err != nil {
				break
			}

			pieces.AddRequested(pieceBlock)
			break
		}
	}
	fmt.Println("Outside for")

	return
}
