package queue

import (
	"fmt"
	"github.com/concurrency-8/parser"
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"testing"
)

func getTorrentFile() parser.TorrentFile {
	torrent, _ := parser.ParseFromFile("../test_torrents/ubuntu.iso.torrent")
	return torrent
}

// TestQueue tests the queue functionality
func TestQueue(t *testing.T) {
	torrent := getTorrentFile()                                                 // get a torrent file
	pieces := math.Floor(float64(torrent.Length / uint64(torrent.PieceLength))) // calculate number of pieces in file

	// tests when making a Queue object
	queue := NewQueue(torrent)
	assert.Equal(t, queue.choked, true, "Choked attribute not true by default")
	assert.Equal(t, queue.length(), 0, "Queue length not zero initially")

	// testing enqueue() function
	indexs := make([]uint32, 0)
	rand.Seed(56)
	epochs := rand.Intn(10)
	totalLength := 0
	for i := 0; i < epochs; i++ {
		index := uint32(rand.Intn(int(pieces))) // generate index between numbe of pieces
		err := queue.enqueue(index)
		assert.Equal(t, err, nil, "Error while enqueue ")
		indexs = append(indexs, index)
		block, err := queue.peek()
		assert.Equal(t, err, nil, "Error encountered while peeking ")
		assert.Equal(t, block.Index, indexs[0], "Front element different in enqueue")
		totalLength = totalLength + int(block.Nblocks)
	}

	assert.Equal(t, queue.length(), totalLength, "Queue length not full after enqueue") // ensure queue is of full length

	// testing dequeue() function
	for i := 0; i < epochs; i++ {
		block, err := queue.peek()
		assert.Equal(t, err, nil, "Error encountered while peeking")

		for j := 0; j < int(block.Nblocks); j++ { // dequeue all blocks corresponding to a piece index
			block, err := queue.peek()
			assert.Equal(t, err, nil, "Error encountered while peeking")
			err = queue.dequeue()
			assert.Equal(t, err, nil, "Error encountered while dequeuing")
			assert.Equal(t, block.Index, indexs[i], "Front element different in dequeue")
		}

	}

	assert.Equal(t, queue.length(), 0, "Queue length not zero finally")

	fmt.Println("PASS")
}
