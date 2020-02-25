package mRedis

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"markV5/mError"
	"markV5/mFace"
	"net"
	"os"
)

type MRedisConfig struct {
	Nickname string
	Network  string
	Host     string
	Port     string

	ReadTimeout    int
	WriteTimeout   int
	ConnectTimeout int

	PoolMaxIdle     int
	PoolMaxActive   int
	PoolIdleTimeout int
	PoolWait        bool
}

func (mrc *MRedisConfig) LoadByJSONFile(path string) mFace.MError {
	file, err := os.Open(path)
	if err != nil {
		return mError.SystemError(err)
	}
	defer file.Close()

	buf := make([]byte, 1024)
	n, err := file.Read(buf)
	if err != nil {
		return mError.SystemError(err)
	}

	err = json.Unmarshal(buf[:n], mrc)
	if err != nil {
		return mError.SystemError(err)
	}

	return mError.Nil
}

func (mrc *MRedisConfig) LoadByTOMLFile(path string) mFace.MError {
	_, err := toml.DecodeFile(path, mrc)
	if err != nil {
		return mError.SystemError(err)
	}
	return mError.Nil
}

func (mrc *MRedisConfig) address() string {
	return net.JoinHostPort(mrc.Host, mrc.Port)
}
