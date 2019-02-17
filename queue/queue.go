package queue

import (
	"fmt"
	"github.com/concurrency-8/parser"
)

// Queue object for storing requested pieces
type Queue struct {
	torrent parser.TorrentFile
	choked  bool
	queue   []PieceBlock
}

// NewQueue returns a fresh pointer to a Queue object
func NewQueue(torrent parser.TorrentFile) (queue *Queue) {
	queue = &Queue{torrent, true, make([]PieceBlock, 0)}
	return
}

// enqueue adds a piece to queue
func (queue *Queue) enqueue(pieceIndex uint32) (err error) {
	nBlocks, err := parser.BlocksPerPiece(queue.torrent, pieceIndex)

	if err != nil {
		return
	}

	for i := 0; i < int(nBlocks); i++ {
		blocklen, err := parser.BlockLen(queue.torrent, pieceIndex, uint32(i))
		if err != nil {
			break
		}

		pieceBlock := PieceBlock{pieceIndex, uint32(i) * parser.BLOCK_LEN, blocklen, nBlocks}
		queue.queue = append(queue.queue, pieceBlock)

	}
	return
}

// dequeue removes first piece block
func (queue *Queue) dequeue() error {
	if queue.length() == 0 {
		return fmt.Errorf("Queue empty : can't dequeue")
	}

	queue.queue = queue.queue[1:]
	return nil
}

// peek returns first pieceblock
func (queue *Queue) peek() (block PieceBlock, err error) {

	if queue.length() == 0 {
		err = fmt.Errorf("Queue empty : can't peek")
	} else {
		block = queue.queue[0]
	}
	return
}

// length returns length of queue
func (queue *Queue) length() int {
	return len(queue.queue)
}
