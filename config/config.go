package config

import (
	"encoding/json"
	"io/ioutil"
)

type Property struct {
	Port    int    `json:"port"`
	Address string `json:"address"`
}

var GlobalConfig *Property

func init() {
	GlobalConfig = &Property{
		Port:    4399,
		Address: "127.0.0.1",
	}
}

func SetupConfig(fileName string) {

	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bytes, GlobalConfig)
	if err != nil {
		panic(err)
	}
}
