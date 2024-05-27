package main

import (
	"burakturkerdev/ftgo/src/common"
	"burakturkerdev/ftgo/src/server/lib"
	"log"
	"net"
	"os"
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
			println("Log => Handshake failed with some client.")
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	c := common.CreateConnection(conn)

	var message common.Message
	c.Read().GetMessage(&message)
	// List dirs operation
	if message == common.CListDirs {

		var path string
		c.GetString(&path)

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
		// End permission checks

		files, err := os.ReadDir(lib.MainConfig.Directory + path)

		if err != nil {
			println("Log => Client is trying to invalid path -> " + err.Error())
		}

		fileinfos := make([]common.FileInfo, len(files))
		//TO-DO file size not working.
		for i, f := range files {
			if !f.IsDir() {
				fileinfos[i] = common.FileInfo{Name: f.Name(), IsDir: f.IsDir(), Size: 0}
			} else {
				stat, err := os.Stat(lib.MainConfig.Directory + string(path) + f.Name())

				if err != nil {
					println("Log => Can't get size of file.")
				}
				fileinfos[i] = common.FileInfo{Name: f.Name(), IsDir: f.IsDir(), Size: stat.Size()}
			}
		}
		c.SendJson(fileinfos)
	}
}
