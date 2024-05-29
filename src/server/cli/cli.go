package main

import (
	"burakturkerdev/ftgo/src/common"
	"burakturkerdev/ftgo/src/server/lib"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

const invalidMsg string = "Invalid message, type ftgo help if you lost."

var resolvers = map[string]common.Resolver{
	"serve":  &ServeResolver{},
	"status": &StatusResolver{},
	"port":   &PortResolver{},
	"dir":    &DirResolver{},
	"perm":   &PermResolver{},
}

// Serve
type ServeResolver struct {
}

func (r ServeResolver) Resolve(head *common.LinkedCommand) {
	o, err := lib.GetDaemonExecCommand().CombinedOutput()

	if err != nil {
		fmt.Println("Error => Daemon can't be started -> " + err.Error())
		return
	}

	fmt.Println("OK => Server started " + string(o))
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

	fmt.Println(build)
}

// Status
type PortResolver struct {
}

func (r PortResolver) Resolve(head *common.LinkedCommand) {
	current := head.Next

	if current == nil {
		fmt.Println(invalidMsg)
		return
	}
	// port command have 3 leafs : add /  rm / list
	add := "add"
	rm := "rm"
	list := "list"

	if current.Command != add && current.Command != rm && current.Command != list {
		fmt.Println(invalidMsg)
		return
	}

	addOrRemove := current.Command == add || current.Command == rm
	// add or remove command takes 1 argument
	if addOrRemove && len(current.Args) != 1 {
		fmt.Println(invalidMsg)
		return
	}

	// add or remove arg should be valid integer
	if addOrRemove {
		if _, err := strconv.Atoi(current.Args[0]); err != nil {
			fmt.Println("Port number should be valid number. Like that: 4040")
			return
		}
	}

	if current.Command == add {
		addingPort := current.Args[0]

		for _, v := range lib.MainConfig.Ports {
			// remove ':' from port for cmp
			existPort := strings.Replace(v, ":", "", -1)

			if addingPort == existPort {
				fmt.Println("Port already exist.")
				return
			}
		}

		lib.MainConfig.Ports = append(lib.MainConfig.Ports, ":"+addingPort)

		err := lib.MainConfig.Save()

		if err != nil {
			fmt.Println("Error => Can't save config file -> " + err.Error())
			return
		}

		fmt.Println(addingPort + " is added.")
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
					fmt.Println("Error => Can't save config file -> " + err.Error())
					return
				}
				fmt.Println(removingPort + " is removed.")
				return
			}
		}
		fmt.Println("Port not found!")
	} else if current.Command == list {
		for i, v := range lib.MainConfig.Ports {
			display := strings.Replace(v, ":", "", -1)
			fmt.Println("Port" + strconv.Itoa(i) + " " + display)
		}
	}
}

// Dir

type DirResolver struct {
}

func (r DirResolver) Resolve(head *common.LinkedCommand) {
	current := head.Next

	if current == nil {
		fmt.Println(invalidMsg)
		return
	}

	//2 leafs, get and set
	get := "get"
	set := "set"

	if current.Command != get && current.Command != set {
		fmt.Println(invalidMsg)
		return
	}

	if current.Command == set && len(current.Args) == 0 {
		fmt.Println("Please give argument for directory. <directory-path>")
		return
	}

	if current.Command == set {
		path := current.Args[0]

		stat, err := os.Stat(path)

		if err != nil {
			fmt.Println("Path is invalid.")
			return
		}

		if !stat.IsDir() {
			fmt.Println("Path should be directory.")
		}

		lib.MainConfig.Directory = path

		err = lib.MainConfig.Save()

		if err != nil {
			fmt.Println("Error => " + err.Error())
			return
		}

		fmt.Println("OK -> Directory set to: " + path)

	} else if current.Command == get {
		fmt.Println("Current directory -> " + lib.MainConfig.Directory)
	}
}

// Perm
type PermResolver struct{}

func (r PermResolver) Resolve(head *common.LinkedCommand) {
	current := head.Next

	if current == nil {
		fmt.Println(invalidMsg)
		return
	}

	// 5 leafs write, read, list, ip, password
	write := "write"
	read := "read"
	list := "list"
	ip := "ip"
	password := "password"

	if current.Command != write && current.Command != read && current.Command != list && current.Command != ip && current.Command != password {
		fmt.Println(invalidMsg)
		return
	}

	if current.Command == write {
		current = current.Next

		if current == nil {
			fmt.Println(invalidMsg)
			return
		}

		// 2 leafs set, get
		set := "set"
		get := "get"

		if current.Command != set && current.Command != get {
			fmt.Println(invalidMsg)
			return
		}

		if current.Command == set {
			// need 1 arg
			if len(current.Args) == 0 {
				fmt.Println("Specify perm for set.")
				return
			}

			perm := lib.WritePerm(current.Args[0])

			if lib.ValidWritePerm(perm) {
				lib.MainConfig.WritePerm = perm
				err := lib.MainConfig.Save()

				if err != nil {
					fmt.Println("Error => " + err.Error())
					return
				}

				fmt.Println("Permission changed.")

			} else {
				fmt.Println("Perm is not exist.")
			}
		} else if current.Command == get {
			fmt.Println(lib.MainConfig.WritePerm)
		}
	} else if current.Command == read {
		current = current.Next

		if current == nil {
			fmt.Println(invalidMsg)
			return
		}

		// 2 leafs set, get
		set := "set"
		get := "get"

		if current.Command != set && current.Command != get {
			fmt.Println(invalidMsg)
			return
		}

		if current.Command == set {
			// need 1 arg
			if len(current.Args) == 0 {
				fmt.Println("Specify perm for set")
				return
			}

			perm := lib.ReadPerm(current.Args[0])

			if lib.ValidReadPerm(perm) {
				lib.MainConfig.ReadPerm = perm
				err := lib.MainConfig.Save()

				if err != nil {
					fmt.Println("Error => " + err.Error())
					return
				}
				fmt.Println("Permission changed.")
			} else {
				fmt.Println("Perm is not exist.")
			}
		} else if current.Command == get {
			fmt.Println(lib.MainConfig.ReadPerm)
		}
	} else if current.Command == list {
		message := "Write Perms:\n"

		for i, v := range lib.WritePermList {
			if i != 0 {
				message = message + ", "
			}
			message = message + string(v)
		}

		message = message + "\n\nRead Perms:\n"

		for i, v := range lib.ReadPermList {
			if i != 0 {
				message = message + ", "
			}
			message = message + string(v)
		}

		fmt.Println(message)
	} else if current.Command == ip {
		current = current.Next

		if current == nil {
			fmt.Println(invalidMsg)
			return
		}

		// has 3 leafs add, rm, list
		add := "add"
		rm := "rm"
		list := list

		if current.Command != add && current.Command != rm && current.Command != list {
			fmt.Println(invalidMsg)
			return
		}

		// add, rm commands take 1 argument.
		if len(current.Args) == 0 && (current.Command == add || current.Command == rm) {
			fmt.Println("Give IP as arg! Like that => 1.1.1.1")
			return
		}

		ipMustValid := func(ip string) {
			result := net.ParseIP(ip)
			if result == nil {
				log.Fatal("This is not a valid IP address.")
			}
		}

		if current.Command == add {
			ip := current.Args[0]
			ipMustValid(ip)

			for _, v := range lib.MainConfig.AllowedIps {
				if ip == v {
					fmt.Println("IP already in allowed IP's.")
					return
				}
			}
			lib.MainConfig.AllowedIps = append(lib.MainConfig.AllowedIps, ip)
			err := lib.MainConfig.Save()
			if err != nil {
				fmt.Println("Error => " + err.Error())
				return
			}

			fmt.Println("OK => IP address added to allowed IP's.")

		} else if current.Command == rm {
			ip := current.Args[0]
			ipMustValid(ip)

			for i, v := range lib.MainConfig.AllowedIps {
				if ip == v {
					hold := lib.MainConfig.AllowedIps[len(lib.MainConfig.AllowedIps)-1]
					lib.MainConfig.AllowedIps[i] = hold
					lib.MainConfig.AllowedIps = lib.MainConfig.AllowedIps[:len(lib.MainConfig.AllowedIps)-1]
					err := lib.MainConfig.Save()
					if err != nil {
						fmt.Println("Error => " + err.Error())
						return
					}
					fmt.Println("OK => IP removed from allowed IP's")
					return
				}
			}
			fmt.Println("IP address not exist in allowed IP's.")
		} else if current.Command == list {
			message := "Allowed IP addresses:\n"

			for i, v := range lib.MainConfig.AllowedIps {
				if i != 0 {
					message = message + ", "
				}
				message = message + string(v)
			}
			fmt.Println(message)
		}
	} else if current.Command == password {
		current = current.Next

		// 1 leaf, set
		set := "set"

		if current.Command != set {
			fmt.Println(invalidMsg)
			return
		}

		setNewPassword := func() {
			fmt.Println("Set password: ")

			password := lib.ReadPassword()

			fmt.Println("Again: ")

			passwordC := lib.ReadPassword()

			if string(password) != string(passwordC) {
				fmt.Println("Passwords doesn't match!")
				return
			}

			pwStr := lib.GenerateHash(password)

			lib.MainConfig.Password = pwStr

			err := lib.MainConfig.Save()

			if err != nil {
				fmt.Println("Error => " + err.Error())
				return
			}

			fmt.Println("OK => Password changed.")
		}

		if lib.MainConfig.Password == "" {
			setNewPassword()
		} else {
			fmt.Println("Enter old password: ")

			password := lib.ReadPassword()

			if lib.ValidateHash([]byte(lib.MainConfig.Password), password) {
				fmt.Println("Old password is not correct.")
				return
			}

			setNewPassword()
		}
	}
}
