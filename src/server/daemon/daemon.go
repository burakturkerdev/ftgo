package main

import (
	"bufio"
	"burakturkerdev/ftgo/src/common"
	"burakturkerdev/ftgo/src/server/lib"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var wg sync.WaitGroup

func startDaemon() {
	if err := lib.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	startServer()
	wg.Wait()
}

func startServer() {
	listeners := make([]net.Listener, len(lib.MainConfig.Ports))

	for i, port := range lib.MainConfig.Ports {
		listener, err := net.Listen("tcp", port)

		if err != nil {
			panic("Error => Network error EC2111")
		}
		listeners[i] = listener
	}

	for _, listener := range listeners {
		wg.Add(1)
		go acceptConnections(listener)
	}
}

func acceptConnections(listener net.Listener) {
	defer wg.Done()
	defer listener.Close()
	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Log => Handshake failed with some client.")
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	c := common.CreateConnection(conn)

	// Gets authentication if needed
	ensureReadAuthentication := func() {
		// Start permission checks
		if lib.MainConfig.ReadPerm == lib.ReadPermPassword {

			c.SendMessage(common.SAuthenticate)

			var password string
			c.Read().IgnoreMessage().GetString(&password)

			if !lib.ValidateHash([]byte(lib.MainConfig.Password), []byte(password)) {
				c.SendMessage(common.SUnAuthorized)
				return
			}
		}
		if lib.MainConfig.ReadPerm == lib.ReadPermIp {
			var allowed bool
			for _, v := range lib.MainConfig.AllowedIps {
				if v == conn.RemoteAddr().String() {
					allowed = true
					break
				}
			}
			if !allowed {
				c.SendMessage(common.SUnAuthorized)
				return
			}
		}

		if lib.MainConfig.ReadPerm == lib.ReadPermNone {
			c.SendMessage(common.SUnAuthorized)
			return
		}
		// End permission checks
	}

	var message common.Message
	c.Read().GetMessage(&message)

	// Set the working path if exists(if client didn't send path string will be "").
	var path string
	c.GetString(&path)
	if !strings.HasPrefix(path, "/") { // any path starting with / would be an absolute path from root, so we just keep the path as it is
		path = filepath.Join(lib.MainConfig.Directory, path)
	}

	// List dirs operation
	switch message {
	case common.CListDirs:
		ensureReadAuthentication()

		files, err := os.ReadDir(path)

		if err != nil {
			fmt.Println("Log => Client is trying to read invalid path -> " + err.Error())
		}

		fileinfos := make([]common.FileInfo, len(files))
		// FIXME: file size not working.
		for i, f := range files {
			if !f.IsDir() {
				fileinfos[i] = common.FileInfo{Name: f.Name(), IsDir: f.IsDir(), Size: 0}
			} else {
				stat, err := os.Stat(path + f.Name())

				if err != nil {
					fmt.Println("Log => Can't get size of file.")
				}
				fileinfos[i] = common.FileInfo{Name: f.Name(), IsDir: f.IsDir(), Size: stat.Size()}
			}
		}
		c.SendJson(fileinfos)
	case common.CDownload:
		ensureReadAuthentication()

		stat, err := os.Stat(path)

		if err != nil {
			c.SendMessageWithData(common.Fail, err.Error())
			return
		}

		if stat.IsDir() {
			c.SendMessageWithData(common.Fail, "This is a directory.")
			return
		}

		file, err := os.Open(path)

		if err != nil {
			c.SendMessageWithData(common.Fail, err.Error())
			return
		}

		defer file.Close()

		buffer := make([]byte, common.ExchangeBufferSize)

		reader := bufio.NewReader(file)

		readLoop := 0

		for {
			_, err := reader.Discard(readLoop * common.ExchangeBufferSize)

			if err != nil {
				c.SendMessageWithData(common.Fail, err.Error())
				return
			}

			readed, err := reader.Read(buffer)

			if err != nil && err.Error() != "EOF" {
				c.SendMessageWithData(common.Fail, err.Error())
				return
			}

			if readed < common.ExchangeBufferSize {
				buffer = buffer[:readed]
				if readed != 0 {
					c.SendData(buffer)
				}
				c.SendMessage(common.Completed)
				return
			}

			c.SendData(buffer)
			readLoop++
		}
	case common.CUpload:
		// Start permission checks
		if lib.MainConfig.WritePerm == lib.WritePermPassword {

			c.SendMessage(common.SAuthenticate)

			var password string
			c.Read().IgnoreMessage().GetString(&password)

			if !lib.ValidateHash([]byte(lib.MainConfig.Password), []byte(password)) {
				c.SendMessage(common.SUnAuthorized)
				return
			}
		}
		if lib.MainConfig.WritePerm == lib.WritePermIp {
			var allowed bool
			for _, v := range lib.MainConfig.AllowedIps {
				if v == conn.RemoteAddr().String() {
					allowed = true
					break
				}
			}
			if !allowed {
				c.SendMessage(common.SUnAuthorized)
				return
			}
		}

		if lib.MainConfig.WritePerm == lib.WritePermPassword {
			c.SendMessage(common.SUnAuthorized)
			return
		}
		// End permission checks

		// Creating dirs (it will not do anything if dirs already exist)
		for i := len(path); i <= 0; i-- {
			if path[i] == '/' {
				os.MkdirAll(path[:i], 0)
				break
			}
		}

		file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0)

		if err != nil {
			c.SendMessageWithData(common.Fail, err.Error())
		}

		//Starting to read buffer

		buffer := make([]byte, common.ExchangeBufferSize)

		var m common.Message

		readStarted := false

		for {
			c.Read().GetMessage(&m)

			if m != common.Success && m != common.Completed {
				c.SendMessageWithData(common.Fail, "message is not valid")
			}

			if m == common.Completed {
				return
			}

			if !readStarted {
				file.Truncate(0)
				readStarted = true
			}

			c.GetData(&buffer)

			_, err = file.Write(buffer)

			if err != nil {
				c.SendMessageWithData(common.Fail, "CRITICIAL "+err.Error())
			}
		}
	default:
		c.SendMessage(common.UnknownMessage)
		fmt.Printf("Log => Client sent unknown message -> %d\n", message)
	}
}
