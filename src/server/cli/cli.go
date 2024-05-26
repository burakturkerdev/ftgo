package main

import (
	"burakturkerdev/ftgo/src/server/lib"
	"fmt"
	"net"
)

var resolvers = map[string]Resolver{
	"serve":  &ServeResolver{},
	"status": &StatusResolver{},
}

type LinkedCommand struct {
	command string
	args    []string
	next    *LinkedCommand
}

type Resolver interface {
	resolve(head *LinkedCommand)
}

// Serve
type ServeResolver struct {
}

func (r ServeResolver) resolve(head *LinkedCommand) {
	o, err := lib.GetDaemonExecCommand().CombinedOutput()

	if err != nil {
		println("Error => Daemon can't be started -> " + err.Error())
		return
	}

	println(string(o))
}

// Status
type StatusResolver struct {
}

func (r StatusResolver) resolve(head *LinkedCommand) {
	var portStatus string
	for _, v := range lib.MainConfig.Ports {
		_, err := net.Dial("tcp", "localhost"+v)
		if err != nil {
			portStatus = fmt.Sprintf(portStatus+"\n"+"Port[%s]= INACTIVE", v)
		} else {

			portStatus = fmt.Sprintf(portStatus+"\n"+"Port[%s]= ACTIVE", v)
		}
	}

	build := "FtGo - Server Status \n" + portStatus + "\n\n" + "Read Permission: " +
		string(lib.MainConfig.ReadPerm) + "\n" + "Write Permission: " + string(lib.MainConfig.WritePerm)

	println(build)
}
