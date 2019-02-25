package torrent

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/concurrency-8/piece"
	"github.com/concurrency-8/queue"
	"github.com/concurrency-8/tracker"
)

type handler func([]byte, net.Conn, *piece.PieceTracker, *queue.Queue, *tracker.ClientStatusReport) error

// Download is a function that handshakes with a peer specified by peer object.
// Concurrently call this function to establish parallel connections to many peers.
func Download(peer tracker.Peer, report *tracker.ClientStatusReport, pieces *piece.PieceTracker) error {
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
	onWholeMessage(conn, msgHandler, pieces, queue, report)
	return err
}
func msgHandler(msg []byte, conn net.Conn, pieces *piece.PieceTracker, queue *queue.Queue, report *tracker.ClientStatusReport) error {

	if (len(msg) == int(uint8(msg[0]))+49) && (bytes.Equal(msg[1:20], []byte("BitTorrent protocol"))) {
		message, err := BuildInterested()
		if err != nil {
			fmt.Println("Error", err.Error())
			return err
		}
		conn.Write(message.Bytes())
	} else {

		_, id, payload := ParseMsg(bytes.NewBuffer(msg))

		if id == 0 {
			ChokeHandler(conn)
		}
		if id == 1 {
			UnchokeHandler(conn, pieces, queue)
		}
		if id == 4 {
			HaveHandler(conn, pieces, queue, payload)
		}
		if id == 5 {
			BitFieldHandler(conn, pieces, queue, payload)
		}
		if id == 7 {
			// TODO pieceHandler
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

			length := int32((buffer.Bytes())[0])
			msgLen = int(length) + 4
		}

		for len(buffer.Bytes()) >= 4 && len(buffer.Bytes()) >= msgLen {
			messageBytes := make([]byte, msgLen)
			binary.Read(buffer, binary.BigEndian, messageBytes)
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
func HaveHandler(conn net.Conn, pieces *piece.PieceTracker, queue *queue.Queue, payload Payload) (pieceIndex uint32, err error) {
	binary.Read(payload["payload"].(*bytes.Buffer), binary.BigEndian, &pieceIndex)
	queueempty := (queue.Length() == 0)
	err = queue.Enqueue(pieceIndex)
	if err != nil {
		return
	}
	if queueempty {
		err = RequestPiece(conn, pieces, queue)
	}
	return
}

// BitFieldHandler handles bitfield protocol
func BitFieldHandler(conn net.Conn, pieces *piece.PieceTracker, queue *queue.Queue, payload Payload) (err error){
	queueempty := (queue.Length() == 0)
	msg := payload["payload"]
	for i, bytevalue := range msg.([]byte) {
		for j := 7; j >= 0; j-- {
			if 1 == bytevalue&1 {
				err = queue.Enqueue(uint32(i*8 + j))
			}
			bytevalue = bytevalue << 1
		}
	}
	if queueempty {
		err = RequestPiece(conn, pieces, queue)
	}

	return
}

// RequestPiece requests a piece
func RequestPiece(conn net.Conn, pieces *piece.PieceTracker, queue *queue.Queue) (err error) {
	if queue.Choked {
		err = fmt.Errorf("Queue is choked")
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
			message, err := BuildRequest(pieceBlock)

			if err != nil {
				break
			}
			_, err = conn.Write(message.Bytes())

			if err != nil {
				fmt.Println(err.Error())
				break
			}
			pieces.AddRequested(pieceBlock)
			break
		}
	}

	return
}
