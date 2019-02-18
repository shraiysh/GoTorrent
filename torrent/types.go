package torrent

import "bytes"

type Payload struct {
	Index int32
	Begin int32
	Block *bytes.Buffer
	Length *bytes.Buffer
}