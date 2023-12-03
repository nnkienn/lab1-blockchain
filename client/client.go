package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/nnkienn/lab1-blockchain/blockchain"
	"github.com/nnkienn/lab1-blockchain/network"
)

var nodes = []string{"127.0.0.1:3001", "127.0.0.1:3002", "127.0.0.1:3003"}
var chainMutex sync.Mutex

func main() {
	chain := blockchain.NewBlockchain()

	for _, node := range nodes {
		go network.StartNode(node, &chainMutex, chain)
	}

	go startClient(&chainMutex, chain)

	select {}
}

func startClient(chainMutex *sync.Mutex, chain *blockchain.Blockchain) {
	fmt.Print("Enter client port: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	clientPort := scanner.Text()

	fmt.Printf("Client started. Listening on port %s\n", clientPort)

	for {
		fmt.Print("Enter command (addtransaction, printchain, hello): ")
		scanner.Scan()
		command := scanner.Text()

		switch command {
		case "addtransaction":
			addTransaction(chainMutex, chain)
		case "printchain":
			printBlockchain(chainMutex, chain)
		case "hello":
			sendHelloRequest(clientPort)
		default:
			fmt.Println("Invalid command")
		}
	}
}

func addTransaction(chainMutex *sync.Mutex, chain *blockchain.Blockchain) {
	fmt.Print("Enter transaction data: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	data := scanner.Text()

	transaction := &blockchain.Transaction{
		Data: []byte(data),
	}

	chainMutex.Lock()
	latestBlock := chain.GetLatestBlock()
	latestBlock.AddTransaction(transaction)
	chainMutex.Unlock()

	fmt.Println("Transaction added.")
}

func printBlockchain(chainMutex *sync.Mutex, chain *blockchain.Blockchain) {
	chainMutex.Lock()
	defer chainMutex.Unlock()

	for _, block := range chain.GetBlocks() {
		blockInfo := fmt.Sprintf("Timestamp: %d\nPrev. hash: %x\n", block.Timestamp, block.PrevBlockHash)
		for _, transaction := range block.Transactions {
			blockInfo += fmt.Sprintf("Transaction: %s\n", string(transaction.Data))
		}
		blockInfo += fmt.Sprintf("Hash: %x\n", block.Hash)

		fmt.Println(blockInfo)
	}
}

func sendHelloRequest(clientPort string) {
	for _, node := range nodes {
		go func(node string) {
			conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", node, clientPort))
			if err != nil {
				fmt.Printf("Error connecting to node at %s: %s\n", node, err.Error())
				return
			}
			defer conn.Close()

			fmt.Fprintf(conn, "hello\n")
			response, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Printf("Error reading response from node at %s: %s\n", node, err.Error())
				return
			}
			fmt.Printf("Response from node at %s: %s", node, response)
		}(node)
	}
}
