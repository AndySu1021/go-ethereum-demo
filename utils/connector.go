package utils

import "github.com/ethereum/go-ethereum/ethclient"

func GetRpcClient(host string) (client *ethclient.Client, err error) {
	client, err = ethclient.Dial(host)
	return
}
