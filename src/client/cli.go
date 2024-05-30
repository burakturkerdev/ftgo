package main

import (
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

		//  TODO
	} else if current.Command == push {

	}
}
