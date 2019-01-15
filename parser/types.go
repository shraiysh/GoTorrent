package parser

import (
			"time"
			"github.com/zeebo/bencode"
)

type FileMetaData struct {
	Path []string `bencode:"path"`
	Length int64 `bencode:"length"`
}

type InfoMetaData struct {
	PieceLength int64 `bencode:"piece length"`
	Pieces []byte `bencode:"pieces"`

	Name string `bencode:"name"`
	Length int64 `bencode:"length"`

	Files bencode.RawMessage `bencode:"files"`
}

type MetaData struct {
	Announce     string `bencode:"announce"`
	AnnounceList [][]string `bencode:"announce-list"`
	Comment      string `bencode:"comment"`
	CreatedBy    string `bencode:"created by"`
	CreatedAt    int64 `bencode:"creation date"`
	Info bencode.RawMessage `bencode:"info"`
}

type File struct {
	Path []string
	Length int64
}

type TorrentFile struct {
	Announce []string
	Comment string
	CreatedBy string
	CreatedAt time.Time
	InfoHash string
	Files []*File
}