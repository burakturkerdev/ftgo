package main

import (
	"burakturkerdev/ftgo/src/common"
	"burakturkerdev/ftgo/src/server/lib"
	"fmt"
	"os"
	"strings"
)

func main() {
	lib.LoadConfig()
	loadResolver().Resolve(loadHeadCommand())
}

func loadResolver() common.Resolver {

	if len(os.Args) <= 1 {
		fmt.Println(invalidMsg)
		os.Exit(0)
	}

	resolver, ok := resolvers[os.Args[1]]

	if !ok {
		fmt.Println(invalidMsg)
		os.Exit(0)
	}
	return resolver
}

func loadHeadCommand() *common.LinkedCommand {
	head := &common.LinkedCommand{}

	var current *common.LinkedCommand = head

	for i := 1; i < len(os.Args); i++ {
		if rune(os.Args[i][0]) == '-' {
			arg := strings.Replace(os.Args[i], `-`, "", -1)
			current.Args = append(current.Args, arg)
		} else {
			if current.Command == "" {
				current.Command = os.Args[i]
			} else {
				current.Next = &common.LinkedCommand{
					Command: os.Args[i],
				}
				current = current.Next
			}
		}
	}
	return head
}
