package tracker

import (
	"encoding/binary"
	"bytes"
)

func BuildConnReq() ([]byte) {
	buffer:= &bytes.Buffer{}
	buffer.Write([]byte{0x00,0x00, 0x04, 0x17, 
						0x27 ,0x10 ,0x19 ,0x80,
						0x00 ,0x00 ,0x00 ,0x00,
						0xa6 ,0xec ,0x6b,0x7d})

	return buffer.Bytes()
}

func ResType(response bytes.Buffer) string {
	action:= binary.BigEndian.Uint32(response.Bytes())
	if (action==0) { return "connect" }
	return "announce" 
}