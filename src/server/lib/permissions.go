package lib

type WritePerm string

type ReadPerm string

const (
	WritePermReadOnly WritePerm = "READONLY"
	WritePermPassword WritePerm = "PASSWORD"
	WritePermEveryone WritePerm = "EVERYONE"
	WritePermIp       WritePerm = "IP"

	ReadPermPassword ReadPerm = "PASSWORD"
	ReadPermIp       ReadPerm = "IP"
	ReadPermEveryone ReadPerm = "EVERYONE"
)
