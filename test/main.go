package main

import (
	"burakturkerdev/ftgo/src/common"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:7373")
	if err != nil {
		panic("connection failed " + err.Error())
	}

	defer conn.Close()

	c := common.CreateConnection(conn)

	c.SendMessageWithData(common.CListDirs, "/test")

	files := []common.FileInfo{}

	var result common.Message

	c.Read().GetMessage(&result).GetJson(&files)

	for _, v := range files {
		println(v.Name + " " + string(rune(v.Size)))
	}
}
