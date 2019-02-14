package queue

import (
	"fmt"
	"github.com/concurrency-8/parser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getTorrentFile() parser.TorrentFile {
	torrent, _ := parser.ParseFromFile("../test_torrents/ubuntu.iso.torrent")
	return torrent
}

// TestQueue tests the queue functionality
func TestQueue(t *testing.T) {
	torrent := getTorrentFile()
	queue := NewQueue(torrent)
	assert.Equal(t, queue.choked, true, "Choked attribute not true by default")
	assert.Equal(t, len(queue.queue), 0, "Queue length not zero initially")
	fmt.Println("PASS")

}
