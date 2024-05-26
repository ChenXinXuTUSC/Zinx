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

		NumWorker: 10,
		NumTaskMx: 1024,

		MaxPacketSize: 2048,
		MaxConn:       1024,
		MaxMsgBuffNum: 32,
	}
	GlobalConfig.Reload()
	log.Info("load config %+v", *GlobalConfig)
}

type Config struct {
	Zinx    zinf.ZinfServer `json:"-"`
	Name    string          `json:"name"`
	Host    string          `json:"host"`
	Port    uint32          `json:"port"`
	Version string          `json:"version"`

	NumWorker uint32 `json:"numworker"`
	NumTaskMx uint32 `json:"numtaskmx"`

	MaxPacketSize uint32
	MaxConn       uint32
	MaxMsgBuffNum uint32
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
