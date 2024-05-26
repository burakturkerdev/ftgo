package main

import (
	"burakturkerdev/ftgo/src/server/lib"
	"os"
	"strings"
)

func main() {
	lib.LoadConfig()
	loadResolver().resolve(loadHeadCommand())
}

func loadResolver() Resolver {

	if len(os.Args) <= 1 {
		println("Type help if you are lost.")
		os.Exit(0)
	}

	resolver, ok := resolvers[os.Args[1]]

	if !ok {
		println(os.Args[1] + " is invalid command. Type help if you are lost.")
		os.Exit(0)
	}
	return resolver
}

func loadHeadCommand() *LinkedCommand {
	head := &LinkedCommand{}

	var current *LinkedCommand = head

	for i := 1; i < len(os.Args); i++ {
		if os.Args[i][0] == '"' {
			arg := strings.Replace(os.Args[i], `"`, "", -1)

			current.args = append(current.args, arg)
		} else {
			current.next = &LinkedCommand{
				command: os.Args[i],
			}
			current = current.next
		}
	}
	return head
}
