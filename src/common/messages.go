package common

import (
	"encoding/binary"
)

type Message uint32

// They should be maximum 1 byte(0-127)
const (
	CListDirs Message = 12
	CDownload Message = 15

	SUnAuthorized Message = 21
	SAuthenticate Message = 22
	Success       Message = 10
	Fail          Message = 22
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
