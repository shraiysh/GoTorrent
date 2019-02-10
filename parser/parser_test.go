package parser

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
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
