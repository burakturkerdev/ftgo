package common

import (
	"encoding/binary"
)

type Message uint32

// They should be maximum 1 byte(0-127)
const (
	// Client side messages
	CListDirs Message = 10
	CDownload Message = 11

	// Server side messages
	SUnAuthorized Message = 12
	SAuthenticate Message = 13

	// General messages
	Success   Message = 14
	Fail      Message = 15
	Completed Message = 16
)

const (
	// 10 MB (DEFAULT)
	ExchangeBufferSize = 1024 * 1024 * 10
)

func messageToBytes(m Message) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(m))
	return buf
}

type FileInfo struct {
	Name  string
	IsDir bool
	Size  int64
}
