package parser

import (
	"time"

	"github.com/zeebo/bencode"
)

//Struct Tags help bencode in identifying which fields to fill.

//FileMetaData contains MetaData about a File.
type FileMetaData struct {
	Path   []string `bencode:"path"`
	Length int64    `bencode:"length"`
}

//InfoMetaData contains MetaData about the torrent.
type InfoMetaData struct {
	PieceLength int64              `bencode:"piece length"`
	Pieces      []byte             `bencode:"pieces"`
	Name        string             `bencode:"name"`
	Length      int64              `bencode:"length"`
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
	Path   []string
	Length int64
}

//TorrentFile contains information about the torrent.
type TorrentFile struct {
	Announce  []string
	Comment   string
	CreatedBy string
	CreatedAt time.Time
	InfoHash  string
	Length    int64
	Files     []*File
}
