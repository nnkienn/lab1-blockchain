package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/your-username/block"
)

var blockchain = &block.BlockChain{}
var mutex = &sync.Mutex{}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Accepted connection from", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		command := scanner.Text()
		fmt.Println("Received command:", command)

		if strings.HasPrefix(command, "ADD_BLOCK") {
			// Example: ADD_BLOCK|transaction1,transaction2,transaction3
			parts := strings.Split(command, "|")
			transactionsData := strings.Split(parts[1], ",")
			var transactions []*block.Transaction
			for _, data := range transactionsData {
				transactions = append(transactions, &block.Transaction{Data: []byte(data)})
			}

			mutex.Lock()
			blockchain.AddBlock(transactions)
			mutex.Unlock()

			fmt.Println("Block added to the blockchain.")
		} else if strings.HasPrefix(command, "ADD_TRANSACTION") {
			// Example: ADD_TRANSACTION|transaction_data
			parts := strings.Split(command, "|")
			transactionData := parts[1]

			mutex.Lock()
			blockchain.AddBlock([]*block.Transaction{{Data: []byte(transactionData)}})
			mutex.Unlock()

			fmt.Println("Transaction added to the blockchain.")
		} else if command == "PRINT_BLOCKCHAIN" {
			mutex.Lock()
			block.PrintBlockchain(blockchain)
			mutex.Unlock()
		} else {
			fmt.Println("Unknown command:", command)
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server started. Listening on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}
