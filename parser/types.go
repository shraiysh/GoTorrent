package parser

import (
	"github.com/zeebo/bencode"
	"math/rand"
	"os"
	"time"
)

//Struct Tags help bencode in identifying which fields to fill.

//FileMetaData contains MetaData about a File.
type FileMetaData struct {
	Path   []string `bencode:"path"`
	Length uint64   `bencode:"length"`
}

//InfoMetaData contains MetaData about the torrent.
type InfoMetaData struct {
	PieceLength uint32             `bencode:"piece length"`
	Piece       []byte             `bencode:"pieces"`
	Name        string             `bencode:"name"`
	Length      uint64             `bencode:"length"`
	Files       bencode.RawMessage `bencode:"files"`
}

//MetaData contains MetaData about the file.
type MetaData struct {
	Announce     string             `bencode:"announce"`
	AnnounceList [][]string         `bencode:"announce-list"`
	Comment      string             `bencode:"comment"`
	CreatedBy    string             `bencode:"created by"`
	CreatedAt    int64              `bencode:"creation date"`
	Info         bencode.RawMessage `bencode:"info"`
}

//File contains length and path of a File in the torrent.
type File struct {
	Path        []string
	Length      uint64
	FilePointer *os.File
}

//TorrentFile contains information about the torrent.
type TorrentFile struct {
	Name 		string
	Announce    []string
	Comment     string
	CreatedBy   string
	CreatedAt   time.Time
	InfoHash    string
	Length      uint64
	Files       []*File
	PieceLength uint32
	Piece       []byte
}

// PieceBlock is struct for a block of a piece
type PieceBlock struct {
	Index   uint32
	Begin   uint32
	Length  uint32
	Nblocks uint32
	Bytes   []byte
}

// Piece Holds a piece
type Piece struct {
	Blocks []PieceBlock
}

// RandomPieceBlock returns a random PieceBlock object from torrent
func RandomPieceBlock(torrent TorrentFile) PieceBlock {
	pieceIndex := rand.Uint32() % uint32(len(torrent.Piece)/20)
	blocksPerPiece, err := BlocksPerPiece(torrent, pieceIndex)
	if err != nil {
		panic(err)
	}
	blockIndex := rand.Uint32() % blocksPerPiece
	blockLength, err := BlockLen(torrent, pieceIndex, blockIndex)
	if err != nil {
		panic(err)
	}
	return PieceBlock{
		Index:   pieceIndex,
		Begin:   blockIndex * BLOCK_LEN,
		Length:  blockLength,
		Nblocks: blocksPerPiece,
	}
}

// GetTorrentFileList gives the list of torrentfiles (these should be alive, in test_torrents)
func GetTorrentFileList() []string {
	return []string{"../test_torrents/ubuntu.iso.torrent", "../test_torrents/big-buck-bunny.torrent"}
}
