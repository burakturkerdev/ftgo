package main

import (
	"bufio"
	"burakturkerdev/ftgo/src/common"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

var resolvers = map[string]common.Resolver{
	"server":  &ServerResolver{},
	"package": &PackageResolver{},
}

type ServerResolver struct{}

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

					var address string

					for n, a := range mainConfig.Servers {
						if n == sv {
							address = a
							break
						}
					}

					if address == "" {
						for i, c := range sv {
							if c == ':' {
								parse := net.ParseIP(sv[:i])

								if parse == nil {
									fmt.Println("This is not a valid ip address and port.")
									return
								}
								address = sv
							}
						}
					}

					if address == "" {
						fmt.Println("This is not a valid ip address and port.")
						return
					}

					dial, err := net.Dial("tcp", address)

					if err != nil {
						fmt.Println(f + " -> error while trying to push file => " + err.Error())
						continue
					}

					defer dial.Close()

					c := common.CreateConnection(dial)

					c.SendMessage(common.CUpload)

					handleAuth(c)

					err = pushFileToServer(f, c)

					if err != nil {
						fmt.Println(f + " -> error while trying to push file => " + err.Error())
					}
				}
			}
		}

	}
}

func pushFileToServer(fp string, c *common.Connection) error {
	file, err := os.OpenFile(fp, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0)

	if err != nil {
		return err
	}

	stat, err := os.Stat(fp)

	if err != nil {
		return err
	}

	reader := bufio.NewReader(file)

	buffer := make([]byte, common.ExchangeBufferSize)

	send := 0

	ch := make(chan int)

	go createProgress(fp, int(stat.Size()/common.ExchangeBufferSize), ch)

	for {
		_, err = reader.Discard(send * common.ExchangeBufferSize)

		if err != nil {
			close(ch)
			return err
		}

		readed, err := reader.Read(buffer)

		if err != nil {
			close(ch)
			return err
		}

		if readed == 0 {
			close(ch)
			c.SendMessage(common.Completed)
			break
		}

		if readed < common.ExchangeBufferSize {
			buffer = buffer[:readed]
		}

		c.SendData(buffer)

		send++
		ch <- send
	}
	close(ch)
	return nil
}

// If authentication needed it will get password from user
func handleAuth(c *common.Connection) {

	var message common.Message

	c.Read().GetMessage(&message)

	if message == common.SAuthenticate {
		pw := common.ReadPassword()

		c.SendString(string(pw))
		c.Read().GetMessage(&message)

		if message == common.Success {
			return
		} else {
			log.Fatal("Can't authenticate with this password.")
		}
	}

}

func createProgress(name string, total int, completed chan int) {
	for range completed {
		msg := name + " -> "
		for range completed {
			msg += "="
		}
		for range total - <-completed {
			msg += "-"
		}
		fmt.Println("\033[2K")
		fmt.Println(msg)
	}
	fmt.Println()
}
