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
	ReadPermNone     ReadPerm = "NONE"
)

var WritePermList = []WritePerm{
	WritePermReadOnly,
	WritePermPassword,
	WritePermEveryone,
	WritePermIp,
}

var ReadPermList = []ReadPerm{
	ReadPermNone,
	ReadPermPassword,
	ReadPermEveryone,
	ReadPermIp,
}

func ValidWritePerm(p WritePerm) bool {
	if p == WritePermEveryone || p == WritePermIp || p == WritePermPassword || p == WritePermReadOnly {
		return true
	}

	return false
}

func ValidReadPerm(p ReadPerm) bool {
	if p == ReadPermEveryone || p == ReadPermIp || p == ReadPermPassword || p == ReadPermNone {
		return true
	}

	return false
}
