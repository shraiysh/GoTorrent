package queue

import (
	"fmt"
	"github.com/concurrency-8/parser"
)

// Queue object for storing requested pieces
type Queue struct {
	torrent parser.TorrentFile
	Choked  bool
	queue   []parser.PieceBlock
}

// NewQueue returns a fresh pointer to a Queue object
func NewQueue(torrent parser.TorrentFile) (queue *Queue) {
	queue = &Queue{torrent, true, make([]parser.PieceBlock, 0)}
	return
}

// Enqueue adds a piece to queue
func (queue *Queue) Enqueue(pieceIndex uint32) (err error) {
	nBlocks, err := parser.BlocksPerPiece(queue.torrent, pieceIndex)

	if err != nil {
		return
	}

	for i := 0; i < int(nBlocks); i++ {
		blocklen, err := parser.BlockLen(queue.torrent, pieceIndex, uint32(i))
		if err != nil {
			break
		}

		pieceBlock := parser.PieceBlock{
			Index:   pieceIndex,
			Begin:   uint32(i) * parser.BLOCK_LEN,
			Length:  blocklen,
			Nblocks: nBlocks,
		}
		queue.queue = append(queue.queue, pieceBlock)

	}
	return
}

// Dequeue removes first piece block
func (queue *Queue) Dequeue() error {
	if queue.Length() == 0 {
		return fmt.Errorf("Queue empty : can't dequeue")
	}

	queue.queue = queue.queue[1:]
	return nil
}

// Peek returns first pieceblock
func (queue *Queue) Peek() (block parser.PieceBlock, err error) {

	if queue.Length() == 0 {
		err = fmt.Errorf("Queue empty : can't peek")
	} else {
		block = queue.queue[0]
	}
	return
}

// Length returns length of queue
func (queue *Queue) Length() int {
	return len(queue.queue)
}
