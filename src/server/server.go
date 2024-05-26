package main

import (
	"burakturkerdev/ftgo/src/common/config"
	"burakturkerdev/ftgo/src/common/connection"
	"burakturkerdev/ftgo/src/common/messages"
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

	c := connection.CreateConnection(conn)

	var message messages.Message

	c.Read().GetMessage(&message)

	// List dirs operation
	if message == messages.CListDirs {
		// Start permission checks
		if mainConfig.ReadPerm == ReadPermPassword {
			c.SendMessage(messages.SAuthenticate)

			var password string

			c.Read().GetString(&password)

			if password != mainConfig.Password {
				c.SendMessage(messages.SUnAuthorized)
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
				c.SendMessage(messages.SUnAuthorized)
				return
			}
		}
		// End permission checks

		var path string

		c.GetString(&path)

		files, err := os.ReadDir(mainConfig.Directory + path)

		if err != nil {
			println("Log => Client is trying to invalid path -> " + err.Error())
		}

		fileinfos := make([]messages.FileInfo, len(files))
		//TO-DO file size not working.
		for i, f := range files {
			if !f.IsDir() {
				fileinfos[i] = messages.FileInfo{Name: f.Name(), IsDir: f.IsDir(), Size: 0}
			} else {
				stat, err := os.Stat(mainConfig.Directory + string(path) + f.Name())

				if err != nil {
					println("Log => Can't get size of file.")
				}
				fileinfos[i] = messages.FileInfo{Name: f.Name(), IsDir: f.IsDir(), Size: stat.Size()}
			}
		}
		c.SendJson(fileinfos)
	}
}
