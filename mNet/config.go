package mNet

import (
	"../mConst"
	"../mFace"
	"encoding/json"
	"log"
	"net"
	"os"
)

func defaultConfig() mFace.MConfig {
	return loadConfigWithFilePath("config.json")
}

func loadConfigWithFilePath(filePath string) mFace.MConfig {
	c := &config{
		filePath:      filePath,
		defaultConfig: mFace.ServerConfig{},
	}
	return c
}

type config struct {
	filePath      string
	defaultConfig mFace.ServerConfig
}

func (c *config) Load() error {
	log.Printf("[%s] load configuration file", mConst.Framework_Name+mConst.Framework_Version)

	file, err := os.Open(c.filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	buf := make([]byte, 1024)
	n, err := file.Read(buf)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(buf[:n], &c.defaultConfig); err != nil {
		return err
	}

	c.Parse()

	log.Printf("[%s] load configuration file success : %+v", c.defaultConfig.Name, c.defaultConfig)

	return nil
}

func (c *config) Reload() error {
	return c.Load()
}

func (c *config) Parse() {
	if c.defaultConfig.Name == "" {
		c.defaultConfig.Name = mConst.Framework_Name + mConst.Framework_Version
	}

	if c.defaultConfig.Network == "" {
		c.defaultConfig.Network = mConst.Network_TCP
	}

	if c.defaultConfig.Host == "" {
		c.defaultConfig.Network = mConst.Default_Host
	}
}

func (c *config) ServerConfig() mFace.ServerConfig {
	return c.defaultConfig
}

func (c *config) ListenServe() bool {
	return c.defaultConfig.Port != ""
}

func (c *config) Address() string {
	return net.JoinHostPort(c.defaultConfig.Host, c.defaultConfig.Port)
}
