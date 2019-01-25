package parser

import (
	"io"
	"io/ioutil"
	"os"
	"time"

	bencode "github.com/zeebo/bencode"
)

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
	if info.Length > 0 {
		files = append(files, &File{
			Path:   []string{info.Name},
			Length: info.Length,
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
			files = append(files, &File{
				Path:   append([]string{info.Name}, f.Path...),
				Length: f.Length,
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
		Announce:  announces,
		Comment:   metadata.Comment,
		CreatedBy: metadata.CreatedBy,
		CreatedAt: time.Unix(metadata.CreatedAt, 0),
		InfoHash:  toSHA1(metadata.Info),
		Length:    Length,
		Files:     files,
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
