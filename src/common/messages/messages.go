package messages

import (
	"bytes"
	"encoding/binary"
)

// Client messages to server all messages exactly should be 10 bytes.
const (
	MessageBufferSize = 10
)

type ClientMessage uint32
type ServerMessage uint32

// They should be maximum 1 byte(0-127)
const (
	CListDirs ClientMessage = 12
	CUpload   ClientMessage = 13
	CDownload ClientMessage = 14

	SUnauthorized ServerMessage = 21
	SAuthenticate ServerMessage = 22
)

func CMessageToBytes(cm ClientMessage) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(cm))
	return buf
}

func SMessageToBytes(sm ServerMessage) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(sm))
	return buf
}

func Trim(b []byte) []byte {
	var path []byte

	nullIndex := bytes.IndexByte(b, 0x00)

	if nullIndex != -1 {
		path = b[:nullIndex]
	} else {
		path = b
	}

	return path
}

type FileInfo struct {
	Name   string
	IsFile bool
	Size   int64
}
