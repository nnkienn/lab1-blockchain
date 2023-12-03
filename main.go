package main

import (
	"fmt"
	"github.com/nnkienn/lab1-blockchain/server/blockchain"
	"github.com/nnkienn/lab1-blockchain/client"
	"github.com/nnkienn/lab1-blockchain/network"
	"sync"
)

var chainMutex sync.Mutex
var nodes = []string{"127.0.0.1:3001", "127.0.0.1:3002", "127.0.0.1:3003"}

func main() {
	chain := block.NewBlockchain()

	for _, node := range nodes {
		go network.StartNode(node, &chainMutex, chain)
	}

	go client.StartClient(&chainMutex, chain)

	select {}
}
