package main

import (
	"os"
	"encoding/json"
	"io/ioutil"
	"os/exec"
)

var ipfsHashes map[string]ipfsHash

type ipfsHash struct {
	Data     string `json:"d"`
	KeyMov   string `json:"k"`
	Combined string `json:"c"`
}

func init() {
	file, err := os.Open("./ipfsHashes.json")
	if os.IsNotExist(err) {
		ipfsHashes = make(map[string]ipfsHash, 100)
		return
	}

	if err != nil {
		panic(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&ipfsHashes)
	if err != nil {
		panic(err)
	}
}

func saveIpfsHashes() error {
	b, _ := json.MarshalIndent(ipfsHashes, "", " ")
	return ioutil.WriteFile("ipfsHashes.json", b, 0644)
}

func ipfsPinPath(path string) (string, error) {
	bin := "ipfs"
	args := []string{"add", "-r", "-p=false", "--nocopy", path}

	ial := exec.Command(bin, args...)
	out, err := ial.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func ipfsAddLink(dirHash string, name string, link string) (string, error) {
	bin := "ipfs"
	args := []string{"object", "patch", "add-link", dirHash, name, link}

	ial := exec.Command(bin, args...)
	out, err := ial.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}