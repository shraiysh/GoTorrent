package parser

import (
	"io"
	"io/ioutil"
	"os"
	"time"
	"github.com/zeebo/bencode"
)

func Parse(reader io.Reader) (*TorrentFile, error) {
	data, err := ioutil.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	metadata := &MetaData{}
	err = bencode.DecodeBytes(data, metadata)
	if err != nil {
		return nil, err
	}

	info := &InfoMetaData{}
	err = bencode.DecodeBytes(metadata.Info, info)
	if err != nil {
		return nil, err
	}

	files := make([]*File, 0)
	// single file context
	if info.Length > 0 {
		files = append(files, &File{
			Path:   []string{info.Name},
			Length: info.Length,
		})
	} else {
		metadataFiles := make([]*FileMetaData, 0)
		err = bencode.DecodeBytes(info.Files, &metadataFiles)
		if err != nil {
			return nil, err
		}

		for _, f := range metadataFiles {
			files = append(files, &File{
				Path:   append([]string{info.Name}, f.Path...),
				Length: f.Length,
			})
		}
	}

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
	
	return &TorrentFile{
		Announce:  announces,
		Comment:   metadata.Comment,
		CreatedBy: metadata.CreatedBy,
		CreatedAt: time.Unix(metadata.CreatedAt, 0),
		InfoHash:  toSHA1(metadata.Info),
		Files:     files,
	}, nil
}

func ParseFromFile(path string) (*TorrentFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return Parse(file)
}
