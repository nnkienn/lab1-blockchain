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
			handleAddBlockCommand(command)
		case strings.HasPrefix(command, "ADD_TRANSACTION"):
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

func handleAddBlockCommand(command string) {
	// Example: ADD_BLOCK|transaction1,transaction2,transaction3
	parts := strings.Split(command, "|")
	transactionsData := strings.Split(parts[1], ",")
	var transactions []*block.Transaction
	for _, data := range transactionsData {
		transactions = append(transactions, &block.Transaction{Data: []byte(data)})
	}

	mutex.Lock()
	bc.AddBlock(transactions)
	mutex.Unlock()

	fmt.Println("Block added to the blockchain.")
}

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

func handleBuildMerkleTreeCommand(conn net.Conn) {
	// Tạo Merkle Tree từ blockchain
	merkleTree := bc.BuildMerkleTree()

	// Gửi Merkle Root cho client
	response := fmt.Sprintf("Merkle root: %x", merkleTree.Root.Data)
	conn.Write([]byte(response + "\n"))

	fmt.Println("Merkle Tree built. Merkle Root sent to the client.")
}

func handleQueryTransactionCommand(command string, conn net.Conn) {
	// Example: QUERY_TRANSACTION|transaction_data
	parts := strings.Split(command, "|")
	transactionData := parts[1]

	// Kiểm tra giao dịch trong Merkle Tree và gửi kết quả về client
	merkleProof := bc.CheckTransactionInMerkleTree(transactionData)
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