package main

import (
	"burakturkerdev/ftgo/src/server/lib"
	"fmt"
	"net"
	"strconv"
	"strings"
)

var invalidMsg string = "Invalid message, type ftgo help if you lost."

var resolvers = map[string]Resolver{
	"serve":  &ServeResolver{},
	"status": &StatusResolver{},
	"port":   &PortResolver{},
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

// Status
type PortResolver struct {
}

func (r PortResolver) resolve(head *LinkedCommand) {
	current := head.next

	if current == nil {
		println(invalidMsg)
		return
	}
	// port command have 3 leafs : add /  rm / list
	add := "add"
	rm := "rm"
	list := "list"

	addOrRemove := current.command == add || current.command == rm
	// add or remove command takes 1 argument
	if addOrRemove && len(current.args) != 1 {
		println(invalidMsg)
		return
	}

	// add or remove arg should be valid integer
	if addOrRemove {
		if _, err := strconv.Atoi(current.args[0]); err != nil {
			println("Port number should be valid number. Like that: 4040")
			return
		}
	}

	if current.command == add {
		addingPort := current.args[0]

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
	} else if current.command == rm {
		removingPort := current.args[0]

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
	} else if current.command == list {
		for i, v := range lib.MainConfig.Ports {
			display := strings.Replace(v, ":", "", -1)
			println("Port" + strconv.Itoa(i) + " " + display)
		}
		return
	}

	println(invalidMsg)
}
