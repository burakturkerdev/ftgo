package lib

import (
	"burakturkerdev/ftgo/src/common"
	"os/exec"
)

var MainConfig *ServerConfig

type ServerConfig struct {
	WritePerm  WritePerm
	ReadPerm   ReadPerm
	Directory  string
	Ports      []string
	AllowedIps []string
	Password   string
	BufferSize int
}

func (c *ServerConfig) SetFieldsToDefault() {
	c.WritePerm = WritePermReadOnly
	c.ReadPerm = ReadPermPassword
	c.Ports = []string{":7373"}
	c.Password = "test"
	c.BufferSize = 2048
	c.Directory = "/home/burak/ftgo/"
	c.AllowedIps = []string{"1.1.1.1"}
}

func LoadConfig() {
	cfg := &ServerConfig{}

	cfg.SetFieldsToDefault()

	common.InitializeConfig[ServerConfig](cfg, ".servercfg")

	MainConfig = cfg
}

func GetDaemonExecCommand() *exec.Cmd {
	return exec.Command("bash", "-c", "nohup ./ftgodaemon > /dev/null 2>&1 &")
}
