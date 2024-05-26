package main

import (
	"burakturkerdev/ftgo/src/common/connection"
	"burakturkerdev/ftgo/src/common/messages"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:7373")
	if err != nil {
		panic("connection failed " + err.Error())
	}

	defer conn.Close()

	c := connection.CreateConnection(conn)

	c.SendMessageWithData(messages.CListDirs, "/test")

	files := []messages.FileInfo{}

	var result messages.Message

	c.Read().GetMessage(&result).GetJson(&files)

	for _, v := range files {
		println(v.Name + " " + string(rune(v.Size)))
	}
}
