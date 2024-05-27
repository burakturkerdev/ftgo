package common

type LinkedCommand struct {
	Command string
	Args    []string
	Next    *LinkedCommand
}

type Resolver interface {
	Resolve(head *LinkedCommand)
}
