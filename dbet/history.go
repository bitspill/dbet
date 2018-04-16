package main

import (
	"os"
	"encoding/json"
	"io/ioutil"
)

var (
	history map[string][]string
)

func init() {

	file, err := os.Open("./history.json")
	if os.IsNotExist(err) {
		history = make(map[string][]string,100)
		return
	}

	if err != nil {
		panic(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&history)
	if err != nil {
		panic(err)
	}
}

func saveHistory() error {
	b, _ := json.MarshalIndent(history, "", " ")
	return ioutil.WriteFile("history.json", b, 0644)
}
