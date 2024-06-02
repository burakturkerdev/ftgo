package main

import (
	"burakturkerdev/ftgo/common"
	"fmt"
	"log"
	"os"
)

const invalidMsg string = "Invalid message, type ftgo help if you lost."

func main() {
	if err := loadConfig(); err != nil {
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
