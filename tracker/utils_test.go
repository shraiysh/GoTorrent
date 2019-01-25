// These are tests for the functions in the file tracker/utils.go
// Run these tests with `go test` in the package directory

package tracker

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func getMockConnectResponse(transactionId uint32, connectionId uint64) bytes.Buffer {
	var mockResponseBuf bytes.Buffer
	var writer *bufio.Writer = bufio.NewWriter(&mockResponseBuf)

	binary.Write(writer, binary.BigEndian, uint32(0)) // action=0 for connect response
	binary.Write(writer, binary.BigEndian, transactionId)
	binary.Write(writer, binary.BigEndian, connectionId)
	writer.Flush()

	return mockResponseBuf
}

func TestBuildConnReq(t *testing.T) {
	fmt.Print("Testing tracker/utils.go : BuildConnReq(): ")
	req := BuildConnReq()
	errorMessage := "Invalid Connection Request for tracker"
	assert.Equal(t, req, []byte{0x00, 0x00, 0x04, 0x17, 0x27, 0x10, 0x19, 0x80, 0x00, 0x00, 0x00, 0x00, 0xa6, 0xec, 0x6b, 0x7d}, errorMessage)
	assert.NotEqual(t, req, []byte{0x01, 0x00, 0x04, 0x17, 0x27, 0x10, 0x19, 0x80, 0x00, 0x00, 0x00, 0x00, 0xa6, 0xec, 0x6b, 0x7d}, errorMessage)
	fmt.Println("Passed")
}

func TestRespType(t *testing.T) {
	fmt.Print("Testing tracker/utils.go : RespType(): ")
	var mockResponseBuf bytes.Buffer
	var writer *bufio.Writer = bufio.NewWriter(&mockResponseBuf)

	// Mock response contains only action - announce
	binary.Write(writer, binary.BigEndian, uint32(1))
	writer.Flush()
	assert.Equal(t, RespType(mockResponseBuf), "announce", "Unable to detect \"announce\" response when action=1")

	// Mock connect response
	mockResponseBuf = getMockConnectResponse(rand.Uint32(), rand.Uint64())
	assert.Equal(t, RespType(mockResponseBuf), "connect", "Unable to detect \"connect\" response when action=0")

	mockResponseBuf.Reset()

	// Mock response has 16 bytes, first 4 bytes show action - announce
	binary.Write(writer, binary.BigEndian, uint32(1)) // 4 bytes written
	for i := 0; i < 3; i++ {                          // Next 12 bytes = 3 uint32
		binary.Write(writer, binary.BigEndian, uint32(rand.Uint32()))
	}
	writer.Flush()
	assert.Equal(t, RespType(mockResponseBuf), "announce", "Unable to detect \"announce\" response when action=1")

	fmt.Println("Passed")
}
