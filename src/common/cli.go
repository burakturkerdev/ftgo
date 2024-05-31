package common

import (
	"log"
	"os"
	"strings"

	"golang.org/x/term"
)

type LinkedCommand struct {
	Command string
	Args    []string
	Next    *LinkedCommand
}

type Resolver interface {
	Resolve(head *LinkedCommand)
}

func LoadHeadCommand(args []string) *LinkedCommand {
	head := &LinkedCommand{}

	var current *LinkedCommand = head

	for i := 1; i < len(args); i++ {
		if rune(args[i][0]) == '-' {
			arg := strings.Replace(args[i], `-`, "", -1)
			current.Args = append(current.Args, arg)
		} else {
			if current.Command == "" {
				current.Command = args[i]
			} else {
				current.Next = &LinkedCommand{
					Command: args[i],
				}
				current = current.Next
			}
		}
	}
	return head
}

func ReadPassword() []byte {
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal("Password can't read.")
	}

	return password
}
