package queue

import "github.com/concurrency-8/parser"

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
func (queue *Queue) enqueue(pieceIndex uint32) {
	nBlocks := parser.BlocksPerPiece(queue.torrent, pieceIndex)

	for i := 0; i < int(nBlocks); i++ {
		pieceBlock := PieceBlock{pieceIndex, uint32(i) * uint32(parser.BLOCK_LEN), uint32(parser.BlockLen(queue.torrent, pieceIndex, uint32(i)))}
		queue.queue = append(queue.queue, pieceBlock)

	}
}

// dequeue removes first piece block
func (queue *Queue) dequeue() {
	queue.queue = queue.queue[1:]
}

// peek returns first pieceblock
func (queue *Queue) peek() PieceBlock {
	return queue.queue[0]
}

// length returns length of queue
func (queue *Queue) length() int {
	return len(queue.queue)
}
