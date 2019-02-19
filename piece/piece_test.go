package piece

import (
	"github.com/concurrency-8/parser"
	"github.com/stretchr/testify/assert"
	"testing"
	// "math/rand"
	// "fmt"
)

func TestNewPieceTracker(t *testing.T) {
	assert := assert.New(t)
	torrent, _ := parser.ParseFromFile("../test_torrents/big-buck-bunny.torrent")
	tracker := NewPieceTracker(torrent)
	for _, i := range tracker.Requested {
		for _, j := range i {
			assert.Equal(j, false)
		}
	}

	for _, i := range tracker.Received {
		for _, j := range i {
			assert.Equal(j, false)
		}
	}
}

func getTorrentBlockTracker() (torrent parser.TorrentFile,
	pieceBlock parser.PieceBlock,
	tracker PieceTracker) {
	torrent, _ = parser.ParseFromFile("../test_torrents/big-buck-bunny.torrent")
	pieceBlock = parser.RandomPieceBlock(torrent)
	tracker = NewPieceTracker(torrent)
	return
}

func TestAddRequested(t *testing.T) {
	_, pieceBlock, tracker := getTorrentBlockTracker()
	tracker.AddRequested(pieceBlock)
	assert.Equal(t, true, tracker.Requested[pieceBlock.Index][pieceBlock.Begin/parser.BLOCK_LEN])
}

func TestAddReceived(t *testing.T) {
	_, pieceBlock, tracker := getTorrentBlockTracker()
	tracker.AddReceived(pieceBlock)
	assert.Equal(t, true, tracker.Received[pieceBlock.Index][pieceBlock.Begin/parser.BLOCK_LEN])
}

func TestIsDone(t *testing.T) {
	_, _, tracker := getTorrentBlockTracker()
	for _, piece := range tracker.Received {
		for blockIndex := range piece {
			piece[blockIndex] = true
		}
	}
	assert.True(t, tracker.IsDone())

}
