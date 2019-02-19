package piece

import (
	"github.com/concurrency-8/parser"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
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

func TestSetRequested(t *testing.T) {
	// _, _, tracker := getTorrentBlockTracker()
	// tracker.setRequested([2][2]bool{{true, false}, {true, false}})
	// assert.Equal(t, true, tracker.Requested[0][0])
	// assert.Equal(t, false, tracker.Requested[0][1])
	// assert.Equal(t, true, tracker.Requested[1][0])
	// assert.Equal(t, false, tracker.Requested[1][1])
}

func TestSetReceived(t *testing.T) {
	// _, _, tracker := getTorrentBlockTracker()
	// tracker.setReceived([][]bool{{true}, {true, false}})
	// assert.Equal(t, tracker.Received, [][]bool{{true}, {true, false}})
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

func TestClone(t *testing.T) {

	// Checking if clone
	src1 := make([][]bool, 5)
	for i := 0; i < 5; i++ {
		src1[i] = make([]bool, 10)
		for j := 0; j < 10; j++ {
			src1[i][j] = rand.Uint32()%2 == 0
		}
	}
	dest1 := clone(src1)

	for i := range src1 {
		for j := range src1[i] {
			assert.Equal(t, src1[i][j], dest1[i][j])
		}
	}

	// Checking if deep clone
	src2 := [][]bool{{true}, {true, false}}
	dest2 := clone(src2)
	src2[1][1] = !src2[1][1]
	assert.NotEqual(t, src2[1][1], dest2[1][1])
}

func TestNeeded(t *testing.T) {
	_, pieceBlock, tracker := getTorrentBlockTracker()

	assert.True(t, tracker.Needed(pieceBlock))

	// Checking if the call for needed copies the received array
	// in requested array when all true

	// The following 2 lines don't work (they don't reflect the changes here)
	// tracker.Requested = [][]bool{{true, true, true}, {true, true, true}}
	// tracker.Received = [][]bool{{true, false, false}, {false, false, true}}

	// Hence, the following 2 lines are used
	tracker.setRequested([][]bool{{true, true, true}, {true, true, true}})
	tracker.setReceived([][]bool{{true, false, false}, {false, false, true}})

	pieceBlock = parser.PieceBlock{
		Index:   1,
		Begin:   200,
		Length:  1000,
		Nblocks: 3,
	}

	tracker.Needed(pieceBlock)

	for i := range tracker.Requested {
		for j := range tracker.Requested[i] {
			assert.Equal(t, tracker.Requested[i][j], tracker.Received[i][j])
		}
	}
}
