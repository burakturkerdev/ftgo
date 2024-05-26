package common

import (
	"encoding/binary"
)

type Message uint32

// They should be maximum 1 byte(0-127)
const (
	CListDirs       Message = 12
	CUploadOverride Message = 13
	CUploadMerge    Message = 14
	CDownload       Message = 15

	SUnAuthorized Message = 21
	SAuthenticate Message = 22
	BufferSize    Message = 24
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
