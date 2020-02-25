package mMysql

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"markV5/mError"
	"markV5/mFace"
	"net/url"
	"os"
	"strings"
)

type MysqlConfig struct {
	NickName string
	User     string
	Password string
	Network  string
	Host     string
	Port     string
	DBName   string
	Options  map[string]string
}

func (mc *MysqlConfig) LoadByJSONFile(path string) mFace.MError {
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

	err = json.Unmarshal(buf[:n], mc)
	if err != nil {
		return mError.SystemError(err)
	}

	return mError.Nil
}

func (mc *MysqlConfig) LoadByTOMLFile(path string) mFace.MError {
	_, err := toml.DecodeFile(path, mc)
	if err != nil {
		return mError.SystemError(err)
	}
	return mError.Nil
}

func (mc *MysqlConfig) dsn() string {
	dataSourceName := fmt.Sprintf("%s:%s@%s(%s:%s)/%s",
		mc.User,
		mc.Password,
		mc.Network,
		mc.Host,
		mc.Port,
		mc.DBName)
	if mc.Options == nil || len(mc.Options) == 0 {
		return dataSourceName
	}

	dataSourceName += "?"
	for key, value := range mc.Options {
		if key == "loc" {
			value = url.QueryEscape(value)
		}
		dataSourceName += fmt.Sprintf("%s=%s&", key, value)
	}
	strings.TrimRight(dataSourceName, "&")

	return dataSourceName
}
