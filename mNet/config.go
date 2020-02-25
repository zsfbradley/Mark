package mNet

import (
	"log"
	"markV5/mError"
	"markV5/mFace"
)

func defaultConfig() mFace.MConfig {
	dc := &config{}
	return dc
}

type config struct {
	name    string
	netWork string
	host    string
	port    string

	rmInCS  uint64
	rmOutCS uint64
	mmInCS  uint64
	mmOutCS uint64
	cmInCS  uint64
	cmOutCS uint64
}

func (c *config) Load() mFace.MError {
	log.Printf("Config Load")
	c.name = "MarkV5"
	c.netWork = "tcp"
	c.host = "0.0.0.0"
	c.port = "8888"
	c.rmInCS = 100
	c.rmOutCS = 100
	c.mmInCS = 100
	c.mmOutCS = 100
	c.cmInCS = 100
	c.cmOutCS = 100
	return mError.Nil
}

func (c *config) Reload() mFace.MError {
	log.Printf("Config Reload")
	c.name = "MarkV5"
	c.netWork = "tcp"
	c.host = "0.0.0.0"
	c.port = "8888"
	c.rmInCS = 100
	c.rmOutCS = 100
	c.mmInCS = 100
	c.mmOutCS = 100
	c.cmInCS = 100
	c.cmOutCS = 100
	return mError.Nil
}

func (c *config) Name() string {
	return c.name
}

func (c *config) NetWork() string {
	return c.netWork
}

func (c *config) Host() string {
	return c.host
}

func (c *config) Port() string {
	return c.port
}

func (c *config) RouteManagerInChanSize() uint64 {
	return c.rmInCS
}

func (c *config) RouteManagerOutChanSize() uint64 {
	return c.rmOutCS
}

func (c *config) MsgManagerInChanSize() uint64 {
	return c.mmInCS
}

func (c *config) MsgManagerOutChanSize() uint64 {
	return c.mmOutCS
}

func (c *config) ConnManagerInChanSize() uint64 {
	return c.cmInCS
}

func (c *config) ConnManagerOutChanSize() uint64 {
	return c.cmOutCS
}
