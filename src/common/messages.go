package common

import (
	"encoding/binary"
)

type Message uint32

// They should be maximum 1 byte(0-127)
const (
	CListDirs Message = 10
	CDownload Message = 11

	SUnAuthorized Message = 12
	SAuthenticate Message = 13
	Success       Message = 14
	Fail          Message = 15
	Blank         Message = 16
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
