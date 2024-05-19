package config

import (
	"encoding/json"
	"os"
	zinf "zinx/interface"
	"zinx/utils/log"
)

func init() {
	// assign default value
	GlobalConfig = &Config{
		Name:    "Zinx",
		Host:    "0.0.0.0",
		Port:    7777,
		Version: "0.4",

		MaxPacketSize: 4096,
		MaxConn:       12000,
	}
	GlobalConfig.Reload()
}

type Config struct {
	Server zinf.ZinfServer `json:"-"`

	Name    string `json:"name"`
	Host    string `json:"host"`
	Port    uint32   `json:"port"`
	Version string `json:"version"`

	MaxPacketSize uint32
	MaxConn       uint32
}

var GlobalConfig *Config

func (cp *Config) Reload() {
	data, readFileErr := os.ReadFile("conf/zinx.json")
	if readFileErr != nil {
		log.Erro(readFileErr.Error())
		panic(readFileErr.Error())
	}

	unmarshalErr := json.Unmarshal(data, GlobalConfig)
	if unmarshalErr != nil {
		log.Erro(unmarshalErr.Error())
		panic(unmarshalErr.Error())
	}
}
