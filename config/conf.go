package config

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

// Conf golbal config setting
var Conf confJSON

// confJSON struct used to parse json
type confJSON struct {
	AccountURL      string        `json:"account_url"`
	URL             string        `json:"url"`
	AddrPath        string        `json:"addr_path"`
	Timeout         time.Duration `json:"timeout"`
	MainPrivateKey  string        `json:"main_private_key"`
	SubPrivateKey   string        `json:"sub_private_key"`
	AirDropAddr     string        `json:"air_drop_addr"`
	DefaultAmount   int64         `json:"default_amount"`
	DefaultGasPrice int64         `json:"default_gas_price"`
	DefaultGasLimit uint64        `json:"default_gas_limit"`
}

// Init initial the config
func Init(filePath string) {
	jsonStr, _ := ioutil.ReadFile(filePath)
	json.Unmarshal(jsonStr, &Conf)
}
