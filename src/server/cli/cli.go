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

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
)

var invalidMsg string = "Invalid message, type ftgo help if you lost."

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

	if current.Command != add && current.Command != rm && current.Command != list {
		println(invalidMsg)
		return
	}

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
	} else if current.Command == list {
		for i, v := range lib.MainConfig.Ports {
			display := strings.Replace(v, ":", "", -1)
			println("Port" + strconv.Itoa(i) + " " + display)
		}
	}
}

// Dir

type DirResolver struct {
}

func (r DirResolver) Resolve(head *common.LinkedCommand) {
	current := head.Next

	if current == nil {
		println(invalidMsg)
		return
	}

	//2 leafs, get and set
	get := "get"
	set := "set"

	if current.Command != get && current.Command != set {
		println(invalidMsg)
		return
	}

	if current.Command == set && len(current.Args) == 0 {
		println("Please give argument for directory. <directory-path>")
		return
	}

	if current.Command == set {
		path := current.Args[0]

		stat, err := os.Stat(path)

		if err != nil {
			println("Path is invalid.")
			return
		}

		if !stat.IsDir() {
			println("Path should be directory.")
		}

		lib.MainConfig.Directory = path

		err = lib.MainConfig.Save()

		if err != nil {
			println("Error => " + err.Error())
			return
		}

		println("OK -> Directory set to: " + path)

	} else if current.Command == get {
		println("Current directory -> " + lib.MainConfig.Directory)
	}
}

// Perm
type PermResolver struct{}

func (r PermResolver) Resolve(head *common.LinkedCommand) {
	current := head.Next

	if current == nil {
		println(invalidMsg)
		return
	}

	// 5 leafs write, read, list, ip, password
	write := "write"
	read := "read"
	list := "list"
	ip := "ip"
	password := "password"

	if current.Command != write && current.Command != read && current.Command != list && current.Command != ip && current.Command != password {
		println(invalidMsg)
		return
	}

	if current.Command == write {
		current = current.Next

		if current == nil {
			println(invalidMsg)
			return
		}

		// 2 leafs set, get
		set := "set"
		get := "get"

		if current.Command != set && current.Command != get {
			println(invalidMsg)
			return
		}

		if current.Command == set {
			// need 1 arg
			if len(current.Args) == 0 {
				println("Specify perm for set.")
				return
			}

			perm := lib.WritePerm(current.Args[0])

			if lib.ValidWritePerm(perm) {
				lib.MainConfig.WritePerm = perm
				err := lib.MainConfig.Save()

				if err != nil {
					println("Error => " + err.Error())
					return
				}

				println("Permission changed.")

			} else {
				println("Perm is not exist.")
			}
		} else if current.Command == get {
			println(lib.MainConfig.WritePerm)
		}
	} else if current.Command == read {
		current = current.Next

		if current == nil {
			println(invalidMsg)
			return
		}

		// 2 leafs set, get
		set := "set"
		get := "get"

		if current.Command != set && current.Command != get {
			println(invalidMsg)
			return
		}

		if current.Command == set {
			// need 1 arg
			if len(current.Args) == 0 {
				println("Specify perm for set")
				return
			}

			perm := lib.ReadPerm(current.Args[0])

			if lib.ValidReadPerm(perm) {
				lib.MainConfig.ReadPerm = perm
				err := lib.MainConfig.Save()

				if err != nil {
					println("Error => " + err.Error())
					return
				}
				println("Permission changed.")
			} else {
				println("Perm is not exist.")
			}
		} else if current.Command == get {
			println(lib.MainConfig.ReadPerm)
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

		println(message)
	} else if current.Command == ip {
		current = current.Next

		if current == nil {
			println(invalidMsg)
			return
		}

		// has 3 leafs add, rm, list
		add := "add"
		rm := "rm"
		list := list

		if current.Command != add && current.Command != rm && current.Command != list {
			println(invalidMsg)
			return
		}

		// add, rm commands take 1 argument.
		if len(current.Args) == 0 && (current.Command == add || current.Command == rm) {
			println("Give IP as arg! Like that => 1.1.1.1")
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
					println("IP already in allowed IP's.")
					return
				}
			}
			lib.MainConfig.AllowedIps = append(lib.MainConfig.AllowedIps, ip)
			err := lib.MainConfig.Save()
			if err != nil {
				println("Error => " + err.Error())
				return
			}

			println("OK => IP address added to allowed IP's.")

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
						println("Error => " + err.Error())
						return
					}
					println("OK => IP removed from allowed IP's")
					return
				}
			}
			println("IP address not exist in allowed IP's.")
		} else if current.Command == list {
			message := "Allowed IP addresses:\n"

			for i, v := range lib.MainConfig.AllowedIps {
				if i != 0 {
					message = message + ", "
				}
				message = message + string(v)
			}
			println(message)
		}
	} else if current.Command == password {
		current = current.Next

		// 1 leaf, set
		set := "set"

		if current.Command != set {
			println(invalidMsg)
			return
		}

		setNewPassword := func() {
			println("Set password: ")

			password, err := terminal.ReadPassword(int(os.Stdin.Fd()))

			if err != nil {
				println("Error reading password:", err)
				return
			}

			println("Again: ")

			passwordC, err := terminal.ReadPassword(int(os.Stdin.Fd()))

			if err != nil {
				println("Error reading password:", err)
				return
			}

			if string(password) != string(passwordC) {
				println("Passwords doesn't match!")
				return
			}

			pwStr := hash(password)

			lib.MainConfig.Password = pwStr

			err = lib.MainConfig.Save()

			if err != nil {
				println("Error => " + err.Error())
				return
			}

			println("OK => Password changed.")
		}

		if lib.MainConfig.Password == "" {
			setNewPassword()
		} else {
			println("Enter old password: ")

			password, err := terminal.ReadPassword(int(os.Stdin.Fd()))

			if err != nil {
				println("Error reading password:", err)
				return
			}

			if validateHash([]byte(lib.MainConfig.Password), password) {
				println("Old password is not correct.")
				return
			}

			setNewPassword()
		}
	}
}

func hash(i []byte) string {
	hash, err := bcrypt.GenerateFromPassword(i, -1)

	if err != nil {
		log.Fatal("SECURITY ERROR / hashing not working.")
	}

	return string(hash)
}

func validateHash(h []byte, i []byte) bool {
	return bcrypt.CompareHashAndPassword(h, i) != nil
}
