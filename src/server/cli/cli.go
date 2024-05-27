package main

import (
	"burakturkerdev/ftgo/src/common"
	"burakturkerdev/ftgo/src/server/lib"
	"fmt"
	"net"
	"strconv"
	"strings"
)

var invalidMsg string = "Invalid message, type ftgo help if you lost."

var resolvers = map[string]common.Resolver{
	"serve":  &ServeResolver{},
	"status": &StatusResolver{},
	"port":   &PortResolver{},
}

// Serve
type ServeResolver struct {
}

func (r ServeResolver) Resolve(head *common.LinkedCommand) {
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

func (r StatusResolver) Resolve(head *common.LinkedCommand) {
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

// Status
type PortResolver struct {
}

func (r PortResolver) Resolve(head *common.LinkedCommand) {
	current := head.Next

	if current == nil {
		println(invalidMsg)
		return
	}
	// port command have 3 leafs : add /  rm / list
	add := "add"
	rm := "rm"
	list := "list"

	addOrRemove := current.Command == add || current.Command == rm
	// add or remove command takes 1 argument
	if addOrRemove && len(current.Args) != 1 {
		println(invalidMsg)
		return
	}

	// add or remove arg should be valid integer
	if addOrRemove {
		if _, err := strconv.Atoi(current.Args[0]); err != nil {
			println("Port number should be valid number. Like that: 4040")
			return
		}
	}

	if current.Command == add {
		addingPort := current.Args[0]

		for _, v := range lib.MainConfig.Ports {
			// remove ':' from port for cmp
			existPort := strings.Replace(v, ":", "", -1)

			if addingPort == existPort {
				println("Port already exist.")
				return
			}
		}

		lib.MainConfig.Ports = append(lib.MainConfig.Ports, ":"+addingPort)

		err := lib.MainConfig.Save()

		if err != nil {
			println("Error => Can't save config file -> " + err.Error())
			return
		}

		println(addingPort + " is added.")
		return
	} else if current.Command == rm {
		removingPort := current.Args[0]

		for i, v := range lib.MainConfig.Ports {
			// remove ':' from port for cmp
			existPort := strings.Replace(v, ":", "", -1)

			if removingPort == existPort {
				lib.MainConfig.Ports[i] = lib.MainConfig.Ports[len(lib.MainConfig.Ports)-1]
				lib.MainConfig.Ports = lib.MainConfig.Ports[:len(lib.MainConfig.Ports)-1]

				err := lib.MainConfig.Save()

				if err != nil {
					println("Error => Can't save config file -> " + err.Error())
					return
				}
				println(removingPort + " is removed.")
				return
			}
		}

		println("Port not found!")
		return
	} else if current.Command == list {
		for i, v := range lib.MainConfig.Ports {
			display := strings.Replace(v, ":", "", -1)
			println("Port" + strconv.Itoa(i) + " " + display)
		}
		return
	}

	println(invalidMsg)
}
