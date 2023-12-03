package main

import (
	"github.com/nnkienn/lab1-blockchain/server/blockchain"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"sync"
)

var nodes = []string{"127.0.0.1:3001", "127.0.0.1:3002", "127.0.0.1:3003"}
var chainMutex sync.Mutex

func main() {
	chain := block.BlockChain{}

	for _, node := range nodes {
		go startServer(node, &chain)
	}

	select {}
}

func startServer(address string, chain *block.BlockChain) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Printf("Node started. Listening on %s\n", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleClient(conn, chain)
	}
}

func handleClient(conn net.Conn, chain *block.BlockChain) {
	fmt.Printf("Client connected: %s\n", conn.RemoteAddr().String())

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		receivedMsg := scanner.Text()
		fmt.Printf("Received message from client: %s\n", receivedMsg)

		if receivedMsg == "printchain" {
			// Send blockchain information to the client
			chainMutex.Lock()
			for _, block := range chain.Blocks {
				blockInfo := fmt.Sprintf("Timestamp: %d\nPrev. hash: %x\n", block.Timestamp, block.PrevBlockHash)
				for _, transaction := range block.Transactions {
					blockInfo += fmt.Sprintf("Transaction: %s\n", string(transaction.Data))
				}
				blockInfo += fmt.Sprintf("Hash: %x\n", block.Hash)

				_, err := fmt.Fprintf(conn, "%s\n", blockInfo)
				if err != nil {
					log.Println(err)
					chainMutex.Unlock()
					return
				}
			}
			chainMutex.Unlock()
		}

		if receivedMsg == "hello" {
			responseMsg := "Hello from node!"
			_, err := fmt.Fprintf(conn, "%s\n", responseMsg)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}
