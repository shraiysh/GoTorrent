package parser

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"regexp"
	"testing"
)

func getTorrentFiles() ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir("../test_torrents")
	return files, err
}

func TestParseFromFile(t *testing.T) {
	files, err := getTorrentFiles()
	assert.Nil(t, err, "opening \"test_torrents\" folder failed.")

	for _, file := range files {
		torfile, err := ParseFromFile("../test_torrents/" + file.Name())
		// check for err!=nil
		assert.Nil(t, err, "Parsing from file failed.")
		// Check for non-empty announce lists
		assert.NotEmpty(t, torfile.Announce, "Empty \"Announce\" list.")
		// There must be atleast one file.
		assert.NotEmpty(t, torfile.Files, "Empty \"File\" list")
		// Length of each file should be positive.
		for _, torsubfile := range torfile.Files {
			assert.True(t, torsubfile.Length > 0, "Negative length for file %s in torrent %s", torsubfile, file.Name())
			assert.NotEmpty(t, torsubfile.Path, "Empty Path for file %s", torsubfile)
		}
		// InfoHash size should be 20 bytes.
		assert.Len(t, torfile.InfoHash, 20, "Corrupt Info Hash file found.")
		// Announce list should consist of valid URLs, i.e. starting with either udp, http or https or wss
		for _, url := range torfile.Announce {
			assert.Regexp(t, regexp.MustCompile("udp://*|http://*|https://*|wss://*"), url, "%s doesn't match any valid tracker format for file %s.", url, file.Name())
		}
		assert.NotEmpty(t, torfile.Length, "Torrent shows empty length.")

	}

}

// TestPieceLen tests PieceLen
func TestPieceLen(t *testing.T) {
	torrent := TorrentFile{}
	torrent.Length = uint64(rand.Intn(100000000))
	torrent.PieceLength = 65536
	lastPieceIndex := uint32(math.Ceil(float64(torrent.Length/uint64(torrent.PieceLength)))) - 1
	lastPieceLen := uint32(torrent.Length % uint64(torrent.PieceLength))

	if lastPieceLen == 0 {
		lastPieceLen = torrent.PieceLength
	}

	for i := 0; i < 2; i++ {
		index := uint32(rand.Intn(int(2 * lastPieceIndex)))
		length, err := PieceLen(torrent, index)
		if index < lastPieceIndex {
			assert.Equal(t, err, nil, "Error not nil")
			assert.Equal(t, length, torrent.PieceLength, "Piece Length not equal")
		} else if index == lastPieceIndex {
			assert.Equal(t, err, nil, "Error not nil")
			assert.Equal(t, length, lastPieceLen, "Piece Length not equal")
		} else {
			assert.NotEqual(t, err, nil, "For large index length still exits")
		}
	}
}

// TestBlocksPerPiece tests BlocksPerPiece
func TestBlocksPerPiece(t *testing.T) {
	torrent := TorrentFile{}
	torrent.Length = uint64(rand.Intn(100000000))
	torrent.PieceLength = 65536
	lastPieceIndex := uint32(math.Ceil(float64(torrent.Length/uint64(torrent.PieceLength)))) - 1

	for i := 0; i < 20; i++ {
		index := uint32(rand.Intn(int(lastPieceIndex)))
		length, err := PieceLen(torrent, index)
		assert.Equal(t, err, nil, "Error not nil")
		blocks, err := BlocksPerPiece(torrent, index)
		assert.Equal(t, err, nil, "Error not nil")
		assert.Equal(t, blocks, uint32(math.Ceil(float64(length)/float64(BLOCK_LEN))), "Block Length not equal")

	}
}

// TestBlockLen tests BlockLen
func TestBlockLen(t *testing.T) {
	torrent := TorrentFile{}
	torrent.Length = uint64(rand.Intn(100000000))
	torrent.PieceLength = 65536
	lastPieceIndex := uint32(math.Ceil(float64(torrent.Length/uint64(torrent.PieceLength)))) - 1
	pieceIndex := uint32(rand.Intn(int(lastPieceIndex)))

	pieceLength, err := PieceLen(torrent, pieceIndex)
	assert.Equal(t, err, nil, "Error not nil")
	lastBlockLength := pieceLength % BLOCK_LEN
	lastBlockIndex := uint32(math.Ceil(float64(pieceLength)/float64(BLOCK_LEN))) - 1

	if lastBlockLength == 0 {
		lastBlockLength = BLOCK_LEN
	}

	for i := 0; i < 20; i++ {
		index := uint32(rand.Intn(int(2 * lastBlockIndex)))
		length, err := BlockLen(torrent, pieceIndex, index)
		if index < lastBlockIndex {
			assert.Equal(t, err, nil, "Error not nil")
			assert.Equal(t, length, BLOCK_LEN, "Block Length not equal")
		} else if index == lastBlockIndex {
			assert.Equal(t, err, nil, "Error not nil")
			assert.Equal(t, length, lastBlockLength, "Block Length not equal")
		} else {
			assert.NotEqual(t, err, nil, "For large index length still exits")
		}
	}

}
