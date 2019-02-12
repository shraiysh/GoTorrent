package torrent

import (
	"encoding/binary"
	"github.com/concurrency-8/parser"
	"math/rand"
)

// GetRandomPiece returns a test piece object with random data.
func GetRandomPiece() (piece parser.Piece) {
	rand.Seed(56)
	piece.Length = rand.Uint32() % 1e6
	piece.Index = rand.Uint32() % 1e6
	piece.Begin = rand.Uint32() % 1e6

	blockLength := rand.Int() % 50
	for i := 0; i < blockLength; i++ {
		binary.Write(&piece.Block, binary.BigEndian, uint8(rand.Int()))
	}

	return
}
