package main

import (
	"burakturkerdev/ftgo/src/common/config"
	"burakturkerdev/ftgo/src/common/messages"
	"bytes"
	"encoding/json"
	"net"
	"os"
	"sync"
)

var mainConfig *ServerConfig

type ServerConfig struct {
	WritePerm  WritePerm
	ReadPerm   ReadPerm
	Directory  string
	Ports      []string
	Password   string
	BufferSize int
}

func (c *ServerConfig) SetFieldsToDefault() {
	c.WritePerm = WritePermReadOnly
	c.ReadPerm = ReadPermPassword
	c.Ports = []string{":7373"}
	c.Password = ""
	c.BufferSize = 2048
	c.Directory = "/home/burak/ftgo/"
}

var wg sync.WaitGroup

func initialize() {
	loadConfig()
	startServer()
	wg.Wait()
}

func loadConfig() {
	cfg := &ServerConfig{}

	cfg.SetFieldsToDefault()

	config.InitializeConfig[ServerConfig](cfg, ".servercfg")

	mainConfig = cfg
}

func startServer() {
	listeners := make([]net.Listener, len(mainConfig.Ports))

	for i, port := range mainConfig.Ports {
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
			println("Log => Handshake failed with some client.")
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	mbuf := make([]byte, 1024)

	_, err := conn.Read(mbuf)

	if err != nil {
		println("Log => Reading from client is failed.")
		return
	}

	// 4. bytes always contain message, first 3 bytes are always 0 bits
	m := uint32(mbuf[3])

	message := messages.ClientMessage(m)

	if message == messages.ListDirs {
		pathBuf := mbuf[4:]

		// Trim zero bits from strings
		var path []byte

		nullIndex := bytes.IndexByte(pathBuf, 0x00)

		if nullIndex != -1 {
			path = pathBuf[:nullIndex]
		} else {
			path = pathBuf
		}
		// Trim end
		pathString := string(path)

		files, err := os.ReadDir(mainConfig.Directory + pathString)

		if err != nil {
			println("Log => Client is trying to invalid path -> " + err.Error())
		}

		//TO-DO PERMISSION CHECKS
		//if perm {}
		fileinfos := make([]messages.FileInfo, len(files))

		//TO-DO file size not working.
		for i, f := range files {
			if !f.IsDir() {
				fileinfos[i] = messages.FileInfo{Name: f.Name(), IsFile: f.IsDir(), Size: 0}
			} else {
				stat, err := os.Stat(mainConfig.Directory + string(path) + f.Name())

				if err != nil {
					println("Log => Can't get size of file.")
				}
				fileinfos[i] = messages.FileInfo{Name: f.Name(), IsFile: f.IsDir(), Size: stat.Size()}
			}
		}

		json, err := json.Marshal(fileinfos)

		if err != nil {
			println("Log => Can't marshal json.")
			return
		}

		conn.Write(json)
	}
}
