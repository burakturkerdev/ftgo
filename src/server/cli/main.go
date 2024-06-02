package main

import (
	"burakturkerdev/ftgo/common"
	"burakturkerdev/ftgo/server/lib"
	"fmt"
	"log"
	"os"
)

func main() {
	if err := lib.LoadConfig(); err != nil {
		log.Fatal(err)
	}
	loadResolver(os.Args).Resolve(common.LoadHeadCommand(os.Args))
}

func loadResolver(args []string) common.Resolver {

	if len(args) <= 1 {
		fmt.Println(invalidMsg)
		os.Exit(0)
	}

	resolver, ok := resolvers[args[1]]

	if !ok {
		fmt.Println(invalidMsg)
		os.Exit(0)
	}
	return resolver
}
