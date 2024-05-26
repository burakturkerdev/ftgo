package lib

import (
	"burakturkerdev/ftgo/src/common"
	"os"
	"os/exec"
	"path/filepath"
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

func (c *ServerConfig) SetFieldsToDefault() error {
	c.WritePerm = WritePermReadOnly
	c.ReadPerm = ReadPermPassword

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	c.Directory = filepath.Join(home, "ftgo")

	c.Ports = []string{":7373"}
	c.AllowedIps = []string{"1.1.1.1"}
	c.Password = "test"
	c.BufferSize = 2048

	return nil
}

func LoadConfig() error {
	cfg := &ServerConfig{}

	if err := cfg.SetFieldsToDefault(); err != nil {
		return err
	}

	common.InitializeConfig[ServerConfig](cfg, ".servercfg")

	MainConfig = cfg

	return nil
}

func GetDaemonExecCommand() *exec.Cmd {
	return exec.Command("bash", "-c", "nohup ./ftgodaemon > /dev/null 2>&1 &")
}
