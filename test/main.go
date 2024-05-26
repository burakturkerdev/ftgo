package main

import (
	"burakturkerdev/ftgo/src/common/messages"
	"encoding/json"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:7373")

	if err != nil {
		panic("connection failed " + err.Error())
	}

	defer conn.Close()
	bytes := messages.CMessageToBytes(messages.CListDirs)
	bytes = append(bytes, []byte("/merhaba")...)
	_, err = conn.Write(bytes)

	if err != nil {
		panic("fail" + err.Error())
	}

	buffer := make([]byte, 1024)

	read, err := conn.Read(buffer)

	if err != nil {
		panic("failhere" + err.Error())
	}

	buffer = buffer[0:read]

	for _, v := range buffer {
		println("%x", v)
	}

	files := []messages.FileInfo{}

	json.Unmarshal(buffer, &files)

	for _, v := range files {
		println(v.Name + " " + string(rune(v.Size)))
	}
}
