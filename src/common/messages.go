package common

import (
	"encoding/binary"
)

// Client messages to server all messages exactly should be 10 bytes.
const (
	MessageBufferSize = 10
)

type Message uint32

// They should be maximum 1 byte(0-127)
const (
	CListDirs Message = 12
	CUpload   Message = 13
	CDownload Message = 14

	SUnAuthorized Message = 21
	SAuthenticate Message = 22
	Success       Message = 10
)

func MessageToBytes(m Message) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(m))
	return buf
}

type FileInfo struct {
	Name  string
	IsDir bool
	Size  int64
}
