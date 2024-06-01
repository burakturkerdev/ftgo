package main

import (
	"bufio"
	"burakturkerdev/ftgo/src/common"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var resolvers = map[string]common.Resolver{
	"server":  &ServerResolver{},
	"package": &PackageResolver{},
	"push":    &PushResolver{},
	"connect": &ConnectResolver{},
	"dir":     &DirResolver{},
}

type ServerResolver struct{}

var pw string = ""

func (r *ServerResolver) Resolve(head *common.LinkedCommand) {
	current := head.Next

	if current == nil {
		fmt.Println(invalidMsg)
		return
	}

	//3 leaf
	add := "add"
	list := "list"
	rm := "rm"

	if current.Command != add && current.Command != list && current.Command != rm {
		fmt.Println(invalidMsg)
		return
	}

	if current.Command == add {
		if len(current.Args) != 2 {
			fmt.Println(invalidMsg)
			return
		}

		name := current.Args[0]
		address := current.Args[1]

		for i, v := range mainConfig.Servers {
			if v == address {
				fmt.Println(v + " already exist with this -> " + i)
				fmt.Println(i + " -> " + v)
				return
			}
			if i == name {
				fmt.Println("Try another name!")
				fmt.Println(i + " already exist with this -> " + v)
				fmt.Println(i + " -> " + v)
				return
			}
		}

		for i, v := range address {
			if v == ':' {
				port := address[i+1:]

				_, err := strconv.Atoi(port)
				if err != nil {
					fmt.Println("Server address is invalid it should look like this: X.X.X.X:PPPP")
					fmt.Println("IPADDRESS:PORT")
					return
				}

				result := net.ParseIP(address[0:i])

				if result == nil {
					log.Fatal("This is not a valid IP address.")
				}
				break
			}
		}

		mainConfig.Servers[name] = address
		err := mainConfig.save()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Server added!")
	} else if current.Command == rm {

		if len(current.Args) != 1 {
			fmt.Println(invalidMsg)
			return
		}

		name := current.Args[0]

		for i := range mainConfig.Servers {
			if i == name {
				delete(mainConfig.Servers, i)
				err := mainConfig.save()
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println(name + " removed from list.")
				return
			}
		}
		fmt.Println(name + " is not exist in server list.")
	} else if current.Command == list {
		msg := "Saved servers:\n"

		for i, v := range mainConfig.Servers {
			msg += i + " => " + v + "\n"
		}

		fmt.Println(msg)
	}
}

// Package resolver
type PackageResolver struct{}

func (p *PackageResolver) Resolve(head *common.LinkedCommand) {
	current := head.Next

	if current == nil {
		fmt.Println(invalidMsg)
		return
	}

	//leafs
	new := "new"
	add := "add"
	rm := "rm"
	list := "list"
	push := "push"

	if current.Command != new && current.Command != add && current.Command != rm && current.Command != list && current.Command != push {
		fmt.Println(invalidMsg)
		return
	}

	if current.Command == new {
		if len(current.Args) != 1 {
			fmt.Println(invalidMsg)
			return
		}

		pname := current.Args[0]

		for _, v := range mainConfig.Packages {
			if v.Name == pname {
				fmt.Println("Name already exist in packages. Use another name.")
				fmt.Println(v.display())
				return
			}
		}

		p := Package{pname, []string{}, time.Now()}

		mainConfig.Packages = append(mainConfig.Packages, &p)

		err := mainConfig.save()

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Package " + pname + " added.")
	} else if current.Command == "add" {
		if len(current.Args) != 2 {
			fmt.Println(invalidMsg)
			return
		}

		pkg := current.Args[0]
		path := current.Args[1]

		for _, v := range mainConfig.Packages {
			if v.Name == pkg {
				for _, k := range v.Files {
					if k == path {
						fmt.Println("File/folder already exist in this package.")
						return
					}
				}

				_, err := os.Stat(path)

				if err != nil {
					fmt.Println("File or folder is not valid. -> " + err.Error())
					return
				}

				v.Files = append(v.Files, path)
				err = mainConfig.save()

				if err != nil {
					fmt.Println(err.Error())
				}

				fmt.Println("Ok => File/folder added to package.")
				return
			}
		}

		fmt.Println("Package " + pkg + " is not exist!")
	} else if current.Command == rm {
		if len(current.Args) != 1 && len(current.Args) != 2 {
			fmt.Println(invalidMsg)
			return
		}

		pkg := current.Args[0]

		for i, v := range mainConfig.Packages {
			if v.Name == pkg {
				if len(current.Args) == 1 {
					for k := i; k < len(mainConfig.Packages)-1; k++ {
						mainConfig.Packages[k] = mainConfig.Packages[k+1]
					}
					mainConfig.Packages = mainConfig.Packages[:len(mainConfig.Packages)-1]
					err := mainConfig.save()

					if err != nil {
						fmt.Println(err.Error())
						return
					}

					fmt.Println("Package " + pkg + " removed.")
					return
				} else if len(current.Args) == 2 {
					f := current.Args[1]

					for l, file := range v.Files {
						if f == file {

							for k := l; k < len(v.Files)-1; k++ {
								v.Files[k] = v.Files[k+1]
							}
							v.Files = v.Files[:len(v.Files)-1]
							err := mainConfig.save()

							if err != nil {
								fmt.Println(err.Error())
								return
							}

							fmt.Println(f + " removed from " + pkg)
							return
						}
					}
					fmt.Println(f + " not exist in " + pkg)
					return
				}
			}
		}
		fmt.Println("Package is not exist.")

	} else if current.Command == list {
		msg := "Current packages:\n"

		for _, v := range mainConfig.Packages {
			msg += v.display() + "\n"
		}

		println(msg)

	} else if current.Command == push {
		if len(current.Args) != 2 {
			fmt.Println(invalidMsg)
			return
		}

		pkgname := current.Args[0]

		sv := current.Args[1]

		for _, v := range mainConfig.Packages {
			if v.Name == pkgname {
				for _, f := range v.Files {

					err := tryUploading(f, sv)

					if err != nil {
						fmt.Println(err.Error())
						return
					}
				}
			}
		}

	}
}

type PushResolver struct{}

func (r *PushResolver) Resolve(head *common.LinkedCommand) {
	current := head

	if len(current.Args) != 2 {
		fmt.Println(invalidMsg)
		return
	}

	file := current.Args[0]
	server := current.Args[1]

	err := tryUploading(file, server)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

}

func tryUploading(file string, serveraddressOrname string) error {
	var address string

	for n, a := range mainConfig.Servers {
		if n == serveraddressOrname {
			address = a
			break
		}
	}

	if address == "" {
		for i, c := range serveraddressOrname {
			if c == ':' {
				parse := net.ParseIP(serveraddressOrname[:i])

				if parse == nil {
					return errors.New("this is not a valid ip address and port")
				}
				address = serveraddressOrname
			}
		}
	}

	if address == "" {
		return errors.New("this is not a valid ip address or server")
	}

	dial, err := net.Dial("tcp", address)

	if err != nil {
		return errors.New(file + " -> error while trying to push file => " + err.Error())
	}

	defer dial.Close()

	c := common.CreateConnection(dial)

	// Exclude home dir for uploading.

	homedir, err := os.UserHomeDir()

	if err == nil {
		c.SendMessageWithString(common.CUpload, strings.Replace(file, homedir, "", -1))
	} else {
		c.SendMessageWithString(common.CUpload, file)
	}

	var m common.Message

	c.Read().GetMessage(&m)

	if m == common.SUnAuthorized {
		return errors.New("access denied")
	}

	if m == common.SAuthenticate {
		m = handleAuth(c)
		if m == common.SUnAuthorized {
			return errors.New("access denied")
		}
	}
	if m != common.Success {
		var s string
		c.GetString(&s)
		return errors.New("error from server -> " + s)
	}
	err = pushFileToServer(file, c)

	if err != nil {
		fmt.Println(file + " -> error while trying to push file => " + err.Error())
	}

	return nil
}

func pushFileToServer(fp string, c *common.Connection) error {
	file, err := os.Open(fp)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	buffer := make([]byte, common.ExchangeBufferSize)

	send := 0

	for {
		readed, err := reader.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}

		if readed == 0 {
			c.SendMessage(common.Completed)
			var m common.Message
			for {
				c.Read().GetMessage(&m)
				if m == common.Completed {
					return nil
				}
			}
		}

		if err == io.EOF {
			c.SendMessage(common.Completed)
			var m common.Message
			for {
				c.Read().GetMessage(&m)
				if m == common.Completed {
					return nil
				}
			}
		}

		// Send data
		c.SendData(buffer[:readed])

		var m common.Message

		// Receive acknowledgment
		c.Read().GetMessage(&m)

		if m != common.Success {
			var d string
			c.GetString(&d)
			return errors.New(d)
		}
		send++
	}
}

type ConnectResolver struct{}

func (r *ConnectResolver) Resolve(head *common.LinkedCommand) {
	current := head

	if len(current.Args) != 1 {
		fmt.Println(invalidMsg)
		return
	}

	server := current.Args[0]

	var address string

	for n, a := range mainConfig.Servers {
		if n == server {
			address = a
			break
		}
	}

	if address == "" {
		for i, c := range server {
			if c == ':' {
				parse := net.ParseIP(server[:i])

				if parse == nil {
					fmt.Println("this is not a valid ip address or port")
					return
				}
				address = server
			}
		}
	}

	if address == "" {
		fmt.Println("this is not a valid ip address or server")
		return
	}

	currentDirectory := "/"

	infos, err := listDirs(currentDirectory, address)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, v := range infos {
		fmt.Println(v.Display())
	}

	cd := "cd"
	pull := "pull"
	exit := "exit"

	parse := func(input string) (string, string) {
		input = strings.Replace(input, "\n", "", -1)
		for i := len(input) - 1; i >= 0; i-- {
			if input[i] == ' ' {
				return input[0:i], input[i+1:]
			}
		}
		return input, ""
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		input, _ := reader.ReadString('\n')

		command, arg := parse(input)

		if command == cd {
			if arg == "" {
				fmt.Println("Commands: cd, pull (fileId/folderId) or exit")
			}

			if arg == ".." {
				if currentDirectory == "/" {
					return
				}

				for i := len(currentDirectory) - 1; i >= 0; i-- {
					if currentDirectory[i] == '/' {
						currentDirectory = currentDirectory[:i]
						break
					}
				}
			} else {
				set := false
				for _, v := range infos {
					if v.IsDir && v.Name == arg {
						if currentDirectory != "/" {
							currentDirectory += "/" + arg
						} else {
							currentDirectory += arg
						}
						set = true
					}
				}
				if !set {
					fmt.Println("This is not a directory.")
					continue
				}
			}
			fmt.Println(currentDirectory)

			infos, err = listDirs(currentDirectory, address)

			if err != nil {
				fmt.Println(err.Error())
				return
			}

			for i := 0; i < 500; i++ {
				fmt.Println("\n")
			}

			for _, v := range infos {
				fmt.Println(v.Display())
			}

		} else if command == pull {
			if arg == "" {
				fmt.Println("Commands: cd, pull (fileId/folderId) or exit")
			}

			valid := false
			for _, v := range infos {
				if v.Name == arg {
					valid = true
					var downloadDirectory string
					if currentDirectory != "/" {
						downloadDirectory = currentDirectory + "/" + arg

					} else {
						downloadDirectory = currentDirectory + arg

					}

					if v.IsDir {

						infos, err = listDirs(downloadDirectory, address)

						if err != nil {
							fmt.Println(arg+" can't be fully downloaded", err.Error())

						}

						recursiveDownload(downloadDirectory, address)
					} else {
						downloadSingle(downloadDirectory, address, v.Size)
					}

					break
				}
			}
			if !valid {
				fmt.Println("This is not valid file or directory.")
			}

		} else if command == exit {
			return
		} else {
			fmt.Println("Commands: cd, pull (fileId/folderId) or exit")
		}

	}

}

type DirResolver struct{}

func (r *DirResolver) Resolve(head *common.LinkedCommand) {
	current := head.Next

	if current == nil {
		fmt.Println(invalidMsg)
		return
	}

	set := "set"
	get := "get"

	if current.Command != get && current.Command != set {
		fmt.Println(invalidMsg)
		return
	}

	if current.Command == set {
		if len(current.Args) != 1 {
			fmt.Println(invalidMsg)
			return
		}

		dir := current.Args[0]
		_, err := os.Stat(dir)

		if err != nil {
			fmt.Println("Invalid path.")
			return
		}

		mainConfig.Downdir = dir
		err = mainConfig.save()

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("Download directory updated.")
	} else if current.Command == get {
		fmt.Println("Download directory -> " + mainConfig.Downdir)
	}

}

func listDirs(path string, server string) ([]common.FileInfo, error) {

	con, err := net.Dial("tcp", server)

	if err != nil {
		return nil, errors.New("error while trying to connect server " + err.Error())

	}

	c := common.CreateConnection(con)

	c.SendMessageWithString(common.CListDirs, path)

	var m common.Message

	c.Read().GetMessage(&m)

	if m == common.SUnAuthorized {
		return nil, errors.New("access denied")
	}

	if m == common.SAuthenticate {
		m = handleAuth(c)
	}
	if m != common.Success {
		var s string
		c.GetString(&s)
		return nil, errors.New("error from server -> " + s)
	}

	var infos []common.FileInfo

	c.GetJson(&infos)

	con.Close()

	return infos, nil
}

func recursiveDownload(path string, server string) error {
	infos, err := listDirs(path, server)
	if err != nil {
		return errors.New("something happened when trying to download this")
	}

	for _, v := range infos {
		newPath := path
		if path == "/" {
			newPath += v.Name
		} else {
			newPath += "/" + v.Name
		}

		if v.IsDir {
			go recursiveDownload(newPath, server)
		} else {
			go downloadSingle(newPath, server, v.Size)
		}
	}
	return nil
}

func downloadSingle(path string, server string, size int64) {
	c, err := net.Dial("tcp", server)

	if err != nil {
		fmt.Println(path+" can't be downloaded.", err.Error())

	}
	defer c.Close()
	con := common.CreateConnection(c)

	con.SendMessageWithString(common.CDownload, path)

	var m common.Message

	con.Read().GetMessage(&m)

	if m == common.SAuthenticate {
		m = handleAuth(con)
	}

	if m != common.Success {
		fmt.Println(path+" can't be downloaded.", err.Error())
		return
	}

	userPath := mainConfig.Downdir + path
	for i := len(userPath) - 1; i >= 0; i-- {
		if userPath[i] == '/' {
			os.MkdirAll(userPath[:i], 0755)
			break
		}
	}

	file, err := os.OpenFile(userPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)

	if err != nil {
		fmt.Println(userPath+" can't be downloaded.", err.Error())
		return
	}
	//Starting to read buffer

	buffer := make([]byte, common.ExchangeBufferSize)

	readStarted := false

	read := 0

	for {
		con.Read().GetMessage(&m)
		if m != common.Completed && !con.EOF {

			if !readStarted {
				file.Truncate(0)
				readStarted = true
			}

			con.GetData(&buffer)
			_, err = file.Write(buffer)

			if err != nil {

				fmt.Println("Can't write to file!")
				return
			}
			con.SendMessage(common.Success)
			read++
		} else {
			fmt.Println(path + " download DONE!")
			con.SendMessage(common.Completed)
			break
		}
	}

}

// If authentication needed it will get password from user
func handleAuth(c *common.Connection) common.Message {
	var m common.Message
	var password string
	if pw == "" {
		fmt.Println("Server wants password:")
		password = string(common.ReadPassword())
		pw = password
	} else {
		password = pw
	}

	c.SendString(password)
	c.Read().GetMessage(&m)
	if m != common.Success {
		log.Fatal("Authentication error.")
	}
	return common.Success
}
