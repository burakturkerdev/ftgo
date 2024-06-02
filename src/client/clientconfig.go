package main

import (
	"burakturkerdev/ftgo/common"
	"os"
)

type ClientConfig struct {
	Servers  map[string]string
	Packages []*Package
	Downdir  string
}

const configPath = ".clientcfg"

var mainConfig *ClientConfig

func setFieldsDefault(c *ClientConfig) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	c.Downdir = home
	c.Packages = []*Package{}
	c.Servers = make(map[string]string)
	return nil
}

func loadConfig() error {
	c := &ClientConfig{}

	err := setFieldsDefault(c)

	if err != nil {
		return err
	}

	common.InitializeConfig[ClientConfig](c, configPath)

	mainConfig = c

	return nil
}

func (c *ClientConfig) save() error {
	return common.SaveConfig(*c, configPath)
}
