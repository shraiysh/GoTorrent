package parser

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"time"
	bencode "github.com/zeebo/bencode"
)

var RESUME = true
// BLOCK_LEN is length  of block
var BLOCK_LEN = uint32(math.Pow(2, 14))

//Parse parses from a stream and returns a pointer to a TorrentFile.
func Parse(reader io.Reader) (TorrentFile, error) {
	data, err := ioutil.ReadAll(reader)
	//return an error if reading fails.
	if err != nil {
		return TorrentFile{}, err
	}

	metadata := &MetaData{}
	err = bencode.DecodeBytes(data, metadata)
	//return an error if decode fails.
	if err != nil {
		return TorrentFile{}, err
	}

	info := &InfoMetaData{}
	err = bencode.DecodeBytes(metadata.Info, info)
	//return an error if further decode fails.
	if err != nil {
		return TorrentFile{}, err
	}
	//This variable refers to the Total torrent size.
	var Length uint64
	files := make([]*File, 0)
	// single file context
	os.Mkdir(info.Name, os.ModePerm)
	if info.Length > 0 {
		var filePointer *os.File
		if !RESUME {
			filePointer, err = os.Create(info.Name + "/" + info.Name)
		}else{
			filePointer, err =  os.OpenFile(info.Name + "/" + info.Name, os.O_APPEND|os.O_WRONLY, 0600)
		}

		if err != nil {
			panic("Unable to create files")
		}
		files = append(files, &File{
			Path:        []string{info.Name},
			Length:      info.Length,
			FilePointer: filePointer,
		})
		Length = info.Length
	} else {
		//multiple files are present.
		metadataFiles := make([]*FileMetaData, 0)
		err = bencode.DecodeBytes(info.Files, &metadataFiles)
		if err != nil {
			return TorrentFile{}, err
		}

		for _, f := range metadataFiles {
			var filePointer *os.File
			if !RESUME {
				filePointer, err = os.Create(info.Name + "/" + f.Path[0])
			}else{
				filePointer, err =  os.OpenFile(info.Name + "/" + f.Path[0], os.O_APPEND|os.O_WRONLY, 0600)
			}
			if err != nil {
				fmt.Println(err)
				panic("Unable to create files ")
			}
			files = append(files, &File{
				Path:        []string{info.Name + "/" + f.Path[0]},
				Length:      f.Length,
				FilePointer: filePointer,
			})
			Length += f.Length
		}
	}

	//announces is the list of trackers.
	announces := make([]string, 0)

	if len(metadata.AnnounceList) > 0 {
		for _, announceItem := range metadata.AnnounceList {
			for _, announce := range announceItem {
				announces = append(announces, announce)
			}
		}
	} else {
		announces = append(announces, metadata.Announce)
	}

	//return the object containing the metadata.
	return TorrentFile{
		Name:		info.Name,
		Announce:    announces,
		Comment:     metadata.Comment,
		CreatedBy:   metadata.CreatedBy,
		CreatedAt:   time.Unix(metadata.CreatedAt, 0),
		InfoHash:    toSHA1(metadata.Info),
		Length:      Length,
		Files:       files,
		PieceLength: info.PieceLength,
		Piece:       info.Piece,
	}, nil
}

//ParseFromFile parses a .torrent file.
func ParseFromFile(path string) (TorrentFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return TorrentFile{}, err
	}
	//Close the file after returning automatically.
	defer file.Close()

	return Parse(file)
}

// PieceLen returns the length of ith piece of file
func PieceLen(torrent TorrentFile, index uint32) (length uint32, err error) {
	totalLength := torrent.Length
	pieceLength := torrent.PieceLength
	lastPieceLen := uint32(totalLength % uint64(pieceLength))
	lastPieceIndex := uint32(math.Ceil(float64(totalLength/uint64(pieceLength)))) - 1

	if lastPieceLen == 0 {
		lastPieceLen = pieceLength
	}

	if lastPieceIndex == index {
		length = lastPieceLen
	} else if lastPieceIndex > index {
		length = pieceLength
	} else {
		err = fmt.Errorf("Piece Index out of range")
	}
	return
}

// BlocksPerPiece returns number of blocks in a ith piece
func BlocksPerPiece(torrent TorrentFile, index uint32) (blocks uint32, err error) {
	pieceLength, err := PieceLen(torrent, index)

	if err != nil {
		return
	}
	blocks = uint32(math.Ceil(float64(pieceLength) / float64(BLOCK_LEN)))
	return
}

// BlockLen calculates length of ith block in jth piece
func BlockLen(torrent TorrentFile, pieceIndex uint32, blockIndex uint32) (length uint32, err error) {
	pieceLength, err := PieceLen(torrent, pieceIndex)

	if err != nil {
		return
	}

	lastBlockLength := pieceLength % BLOCK_LEN
	lastBlockIndex := uint32(math.Ceil(float64(pieceLength)/float64(BLOCK_LEN))) - 1

	if lastBlockLength == 0 {
		lastBlockLength = BLOCK_LEN
	}

	if lastBlockIndex == blockIndex {
		length = lastBlockLength
	} else if lastBlockIndex > blockIndex {
		length = BLOCK_LEN
	} else {
		err = fmt.Errorf("Block Index out of range")
	}
	return
}
