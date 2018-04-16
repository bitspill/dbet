package main

import (
	"github.com/bitspill/flojson"
	"errors"
)

var (
	id     int64
	user   string
	pass   string
	server string
)

func init() {
	id = 0
	user = config.FloConfiguration.RpcUser
	pass = config.FloConfiguration.RpcPass
	server = config.FloConfiguration.RpcAddress
}

func signMessage(address, message string) (string, error) {
	cmd, err := flojson.NewSignMessageCmd(id, address, message)
	if err != nil {
		return "", err
	}

	reply, err := sendRPC(cmd)
	if err != nil {
		return "", err
	}

	if signature, ok := reply.Result.(string); ok {
		return signature, nil
	}

	return "", errors.New("unexpected rpc error")
}

func sendToAddress(address string, amount float64, floData string) (string, error) {
	satoshi := int64(1e8 * amount)
	cmd, err := flojson.NewSendToAddressCmd(id, address, satoshi, "", "", floData)
	if err != nil {
		return "", err
	}

	reply, err := sendRPC(cmd)
	if err != nil {
		return "", err
	}
	if reply.Error != nil {
		return "", reply.Error
	}
	return reply.Result.(string), nil
}

func sendRPC(cmd flojson.Cmd) (flojson.Reply, error) {
	reply, err := flojson.RpcSend(user, pass, server, cmd)
	if err != nil {
		return reply, err
	}
	if reply.Error != nil {
		return reply, reply.Error
	}
	return reply, nil
}


