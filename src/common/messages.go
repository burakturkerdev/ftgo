package common

import (
	"encoding/binary"
	"strconv"
)

type Message uint32

// They should be maximum 1 byte(0-127)
const (
	// Client side messages
	CListDirs Message = 10
	CDownload Message = 11
	CUpload   Message = 12

	// Server side messages
	SUnAuthorized Message = 13
	SAuthenticate Message = 14

	// General messages
	Success        Message = 15
	Fail           Message = 16
	Completed      Message = 17
	UnknownMessage Message = 18
)

const (
	ExchangeBufferSize = 1024 * 10
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

func (f *FileInfo) Display() string {
	msg := f.Name
	if f.IsDir {
		msg += "/"
	} else {
		msg += " " + strconv.Itoa(int(f.Size/(1024*1024))) + " MB"
	}
	return msg
}
