package messages

import "encoding/binary"

// Client messages to server all messages exactly should be 10 bytes.
const (
	MessageBufferSize = 10
)

type ClientMessage uint32

// They should be maximum 1 byte(0-127)
const (
	ListDirs ClientMessage = 12
	Upload   ClientMessage = 13
	Download ClientMessage = 14
)

func MessageToBytes(m ClientMessage) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(m))
	return buf
}

type FileInfo struct {
	Name   string
	IsFile bool
	Size   int64
}
