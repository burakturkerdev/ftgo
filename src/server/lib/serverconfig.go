package lib

import (
	"burakturkerdev/ftgo/src/common"
	"os"
	"os/exec"
	"path/filepath"
)

var cfgPath string = ".servercfg"
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

func (c *ServerConfig) Save() error {
	return common.SaveConfig(*c, cfgPath)
}

func (c *ServerConfig) SetFieldsToDefault() error {
	c.WritePerm = WritePermEveryone
	c.ReadPerm = ReadPermEveryone

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	c.Directory = filepath.Join(home, "ftgo")

	c.Ports = []string{":7373"}
	c.AllowedIps = []string{"1.1.1.1"}
	c.BufferSize = 2048

	return nil
}

func LoadConfig() error {
	cfg := &ServerConfig{}

	if err := cfg.SetFieldsToDefault(); err != nil {
		return err
	}

	common.InitializeConfig[ServerConfig](cfg, cfgPath)

	MainConfig = cfg

	return nil
}

// It will be implemented as os - specific.
func GetDaemonExecCommand() *exec.Cmd {
	return exec.Command("bash", "-c", "nohup ./ftgodaemon > /dev/null 2>&1 &")
}
