package main

import (
	"os"
	"encoding/json"
)

var (
	config Configuration
)

type Configuration struct {
	DatabaseConfiguration `json:"databaseConfiguration"`
	FloConfiguration      `json:"floConfiguration"`
}

type FloConfiguration struct {
	FloAddress string  `json:"floAddress"`
	RpcAddress string  `json:"rpcAddress"`
	RpcUser    string  `json:"rpcUser"`
	RpcPass    string  `json:"rpcPass"`
	TxFeePerKb float64 `json:"txFeePerKb"`
}

type DatabaseConfiguration struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Net      string `json:"net"`
	Address  string `json:"address"`
	Name     string `json:"name"`
}

func init() {
	file, err := os.Open("./conf.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}
}
