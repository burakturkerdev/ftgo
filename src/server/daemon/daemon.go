package main

import (
	"bufio"
	"burakturkerdev/ftgo/src/common"
	"burakturkerdev/ftgo/src/server/lib"
	"fmt"
	"io"
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

	var absolutePath string

	var clientPath string
	c.GetString(&clientPath)

	if !strings.HasPrefix(absolutePath, "/") {
		absolutePath = filepath.Join(lib.MainConfig.Directory, clientPath)
	}

	// List dirs operation
	switch message {
	case common.CListDirs:
		ensureReadAuthentication()
		files, err := os.ReadDir(absolutePath)

		if err != nil {
			fmt.Println("Log => Client is trying to read invalid path -> " + err.Error())
		}

		fileinfos := make([]common.FileInfo, len(files))

		for i, f := range files {
			if !f.IsDir() {
				var size int64 = 0

				info, err := f.Info()

				if err == nil {
					size = info.Size()
				}
				fileinfos[i] = common.FileInfo{Name: f.Name(), IsDir: f.IsDir(), Size: size}
			} else {
				//stat, err := os.Stat(absolutePath + f.Name())

				if err != nil {
					fmt.Println("Log => Can't get size of file.")
				}
				fileinfos[i] = common.FileInfo{Name: f.Name(), IsDir: f.IsDir(), Size: 0}
			}
		}
		c.SendJson(fileinfos)
	case common.CDownload:
		ensureReadAuthentication()
		//We are ready to sending data
		c.SendMessage(common.Success)

		stat, err := os.Stat(absolutePath)

		if err != nil {
			c.SendMessageWithString(common.Fail, err.Error())
			return
		}

		if stat.IsDir() {
			c.SendMessageWithString(common.Fail, "This is a directory.")
			return
		}

		file, err := os.Open(absolutePath)

		if err != nil {
			c.SendMessageWithString(common.Fail, err.Error())
			return
		}

		defer file.Close()

		buffer := make([]byte, common.ExchangeBufferSize)

		reader := bufio.NewReader(file)

		readLoop := 0

		for {

			readed, err := reader.Read(buffer)

			if err != nil && err != io.EOF {
				c.SendMessageWithString(common.Fail, err.Error())
				return
			}

			if readed < common.ExchangeBufferSize {
				buffer = buffer[:readed]
				if readed != 0 {
					c.SendData(buffer)
				}

				c.Read().GetMessage(&message)

				if message != common.Success {
					c.SendMessageWithString(common.Fail, "Corrupted data recieved.")
					return
				}
				for {
					c.SendMessage(common.Completed)
					c.Read().GetMessage(&message)
					if message == common.Completed {
						return
					}
				}

			}

			c.SendData(buffer)

			c.Read().GetMessage(&message)

			if message != common.Success {
				c.SendMessageWithString(common.Fail, "Corrupted data recieved.")
				return
			}
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

		if lib.MainConfig.WritePerm == lib.WritePermReadOnly {
			c.SendMessage(common.SUnAuthorized)
			return
		}
		// End permission checks

		// Is file name provided?
		if clientPath == "" {
			c.SendMessageWithString(common.Fail, "specify file name")
			return
		}

		// Creating dirs (it will not do anything if dirs already exist)
		for i := len(absolutePath) - 1; i >= 0; i-- {
			if absolutePath[i] == '/' {
				os.MkdirAll(absolutePath[:i], 0755)
				break
			}
		}

		file, err := os.OpenFile(absolutePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)

		if err != nil {
			c.SendMessageWithString(common.Fail, err.Error())
			return
		}
		// We are ready for data!
		c.SendMessage(common.Success)

		//Starting to read buffer

		buffer := make([]byte, common.ExchangeBufferSize)

		var m common.Message

		readStarted := false

		for {
			c.Read().GetMessage(&m)
			if m != common.Completed && !c.EOF {
				if !readStarted {
					file.Truncate(0)
					readStarted = true
				}

				c.GetData(&buffer)
				_, err = file.Write(buffer)

				if err != nil {
					c.SendMessageWithString(common.Fail, "CRITICAL "+err.Error())
					return
				}
				c.SendMessage(common.Success)
			} else {
				c.SendMessage(common.Completed)
				break
			}
		}
	default:
		c.SendMessage(common.UnknownMessage)
		fmt.Printf("Log => Client sent unknown message -> %d\n", message)
	}
}
