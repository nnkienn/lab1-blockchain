package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/nnkienn/lab1-blockchain/blockchain"
)

var bc = &block.BlockChain{}
var mutex = &sync.Mutex{}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Accepted connection from", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		command := scanner.Text()
		fmt.Println("Received command:", command)

		switch {
		case strings.HasPrefix(command, "ADD_BLOCK"):
			handleAddTransactionCommand(command)
		case command == "PRINT_BLOCKCHAIN":
			handlePrintBlockchainCommand()
		case command == "BUILD_MERKLE_TREE":
			handleBuildMerkleTreeCommand(conn)
		case strings.HasPrefix(command, "QUERY_TRANSACTION"):
			handleQueryTransactionCommand(command, conn)
		default:
			fmt.Println("Unknown command:", command)
		}
	}
}
// server.go

// ...

func handleBuildMerkleTreeCommand(conn net.Conn) {
	// Tạo Merkle Tree từ blockchain
	merkleTree := bc.BuildMerkleTree()

	// Gửi Merkle Root cho client
	response := fmt.Sprintf("%x", merkleTree.Root.Data)
	conn.Write([]byte(response + "\n"))

	fmt.Println("Merkle Tree built. Merkle Root sent to the client.")
}

// ...



func handleAddTransactionCommand(command string) {
	// Example: ADD_TRANSACTION|transaction_data
	parts := strings.Split(command, "|")
	transactionData := parts[1]

	mutex.Lock()
	bc.AddBlock([]*block.Transaction{{Data: []byte(transactionData)}})
	mutex.Unlock()

	fmt.Println("Transaction added to the blockchain.")
}

func handlePrintBlockchainCommand() {
	mutex.Lock()
	block.PrintBlockchain(bc)
	mutex.Unlock()
}



func handleQueryTransactionCommand(command string, conn net.Conn) {
	// Example: QUERY_TRANSACTION|transaction_data
	parts := strings.Split(command, "|")
	transactionData := parts[1]

	mutex.Lock()
	merkleProof := bc.CheckTransactionInMerkleTree(transactionData)
	mutex.Unlock()

	response := fmt.Sprintf("Transaction verification result: %t", merkleProof)
	conn.Write([]byte(response + "\n"))
}

func main() {
	// Sử dụng flag để nhận giá trị cổng từ dòng lệnh
	port := flag.Int("port", 8080, "Server port")
	flag.Parse()

	serverAddress := fmt.Sprintf("localhost:%d", *port)

	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Server started. Listening on :%d\n", *port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}
