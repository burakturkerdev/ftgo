package main

import (
	"burakturkerdev/ftgo/src/common"
	"fmt"
	"log"
	"net"
	"strconv"
)

var resolvers = map[string]common.Resolver{
	"server": &ServerResolver{},
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
