package piece

import (
	"fmt"
	"math"
	"sync"

	"github.com/concurrency-8/parser"
)

// PieceTracker stores flags for blocks of pieces requested and received
// Requested[i][j] = true => jth block of ith piece has been requested
type PieceTracker struct {
	Torrent   parser.TorrentFile
	Requested [][]bool
	Received  [][]bool
	Lock      sync.Mutex
}

// NewPieceTracker returns a new PieceTracker object for the torrent
func NewPieceTracker(torrent parser.TorrentFile) (tracker *PieceTracker) {
	tracker = new(PieceTracker)
	tracker.Torrent = torrent
	numPieces := uint32(len(torrent.Piece) / 20)
	for i := uint32(0); i < numPieces; i++ {
		blocksPerPiece, _ := parser.BlocksPerPiece(torrent, i)
		tracker.Requested = append(tracker.Requested, make([]bool, blocksPerPiece))
		tracker.Received = append(tracker.Received, make([]bool, blocksPerPiece))
	}

	return
}

// AddRequested flags the request value of a block in a piece
// Invoked while requesting the block of a piece
func (tracker *PieceTracker) AddRequested(block parser.PieceBlock) {
	index := block.Begin / parser.BLOCK_LEN
	tracker.Requested[block.Index][index] = true
}

// AddReceived flags the received value of a block in a piece
// Invoked when a block is received
func (tracker *PieceTracker) AddReceived(block parser.PieceBlock) {
	index := block.Begin / parser.BLOCK_LEN
	tracker.Received[block.Index][index] = true
}

// Needed checks if we want a block. If we have already requested all,
// we reset requested to be equal to received and request the remaining pieces
// Not putting locks here, the caller must make sure that this runs at once by a single thread
func (tracker *PieceTracker) Needed(block parser.PieceBlock) bool {

	// Check if all have been requested...
	allRequested := true
	for _, i := range tracker.Requested {
		if !allRequested {
			break
		}
		for _, j := range i {
			allRequested = allRequested && j
		}
	}

	// If yes, copy received into request...
	if allRequested {
		tracker.Requested = clone(tracker.Received)
	}

	return !tracker.Requested[block.Index][block.Begin/parser.BLOCK_LEN]
}

// Deep clones 2-D bool array
func clone(array [][]bool) (result [][]bool) {
	for _, i := range array {
		temp := make([]bool, len(i))
		for index, j := range i {
			temp[index] = j
		}
		result = append(result, temp)
	}
	return
}

// PieceIsDone tells if the pieceIndex piece has been downloaded successfully
func (tracker *PieceTracker) PieceIsDone(pieceIndex uint32) (result bool) {
	tracker.Lock.Lock()
	result = true
	for _, i := range tracker.Received[pieceIndex] {
		result = result && i
	}
	tracker.Lock.Unlock()
	return
}

// IsDone tells if the torrent file has been successfully received
func (tracker *PieceTracker) IsDone() (result bool) {
	tracker.Lock.Lock()
	result = true
	for _, i := range tracker.Received {
		for _, j := range i {
			result = result && j
		}
	}
	tracker.Lock.Unlock()
	return
}

// PrintPercentageDone prints the percentage of download completed on the screen
func (tracker *PieceTracker) PrintPercentageDone() (percent int) {
	downloaded, total := 0.0, 0
	for _, i := range tracker.Received {
		for _, j := range i {
			total++
			if j {
				downloaded++
			}
		}
	}
	percent = int(math.Round(float64(downloaded*100) / float64(total)))
	// fmt.Print("progress:", percent, "\r")
	return
}

// Reset the piece - Called when invalid SHA
func (tracker *PieceTracker) Reset(index uint32) {
	tracker.Lock.Lock()
	for i := range tracker.Requested[index] {
		tracker.Requested[index][i] = false
		tracker.Received[index][i] = false
	}
	tracker.Lock.Unlock()
}

// Fill is used to revive the piecetracker while resuming the torrent
func (tracker *PieceTracker) Fill(index uint32) {
	for i := range tracker.Requested[index] {
		tracker.Requested[index][i] = true
		tracker.Received[index][i] = true
	}
}

// PrintLeft prints left
func (tracker *PieceTracker) PrintLeft() {
	for i := range tracker.Received {
		for j := range tracker.Received[i] {
			if !tracker.Received[i][j] {
				fmt.Print("[", i, "][", j, "]\t")
			}
		}
	}
	fmt.Println()
}
