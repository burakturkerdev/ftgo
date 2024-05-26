package main

import (
	"burakturkerdev/ftgo/src/common/config"
	"burakturkerdev/ftgo/src/common/messages"
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
	AllowedIps []string
	Password   string
	BufferSize int
}

func (c *ServerConfig) SetFieldsToDefault() {
	c.WritePerm = WritePermReadOnly
	c.ReadPerm = ReadPermPassword
	c.Ports = []string{":7373"}
	c.Password = "test"
	c.BufferSize = 2048
	c.Directory = "/home/burak/ftgo/"
	c.AllowedIps = []string{"1.1.1.1"}
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

	// 4. bytes always contain message, first 3 bytes are always 0 bits they just lead to message.
	m := uint32(mbuf[3])

	message := messages.ClientMessage(m)

	// List dirs operation
	if message == messages.CListDirs {
		//For listing path, message always leading to path. So after 4 bytes, other bytes contain query path.
		pathBuf := mbuf[4:]

		pathString := string(messages.Trim(pathBuf))

		files, err := os.ReadDir(mainConfig.Directory + pathString)

		if err != nil {
			println("Log => Client is trying to invalid path -> " + err.Error())
		}

		// Permission checks
		if mainConfig.ReadPerm == ReadPermPassword {
			conn.Write(messages.SMessageToBytes(messages.SAuthenticate))

			password := make([]byte, 1024)

			conn.Read(password)

			if string(messages.Trim(password)) != mainConfig.Password {
				conn.Write(messages.SMessageToBytes(messages.SUnauthorized))
				return
			}
		}
		if mainConfig.ReadPerm == ReadPermIp {
			var allowed bool
			for _, v := range mainConfig.AllowedIps {
				if v == conn.RemoteAddr().String() {
					allowed = true
					break
				}
			}
			if !allowed {
				conn.Write(messages.SMessageToBytes(messages.SUnauthorized))
				return
			}
		}
		// Permission checks

		fileinfos := make([]messages.FileInfo, len(files))

		//TO-DO file size not working.
		for i, f := range files {
			if !f.IsDir() {
				fileinfos[i] = messages.FileInfo{Name: f.Name(), IsFile: f.IsDir(), Size: 0}
			} else {
				stat, err := os.Stat(mainConfig.Directory + string(pathString) + f.Name())

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
