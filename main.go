package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/nnkienn/lab1-blockchain/server/blockchain"
)

var nodes = []string{"127.0.0.1:3001", "127.0.0.1:3002", "127.0.0.1:3003"}
var chainMutex sync.Mutex

func main() {
	chain := blockchain.NewBlockchain() // Sửa đổi tên package

	for _, node := range nodes {
		go startServer(node, chain)
	}

	// Chờ người dùng nhập lệnh từ terminal
	handleUserInput(chain)
}

func startServer(address string, chain *blockchain.Blockchain) {
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

func handleUserInput(chain *blockchain.Blockchain) {
	for {
		fmt.Print("Enter command (addtransaction, printchain, hello, exit): ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		command := scanner.Text()

		switch command {
		case "addtransaction":
			fmt.Print("Enter transaction data: ")
			scanner.Scan()
			data := scanner.Text()
			handleAddTransaction(data, chain)
		case "printchain":
			handlePrintChain(chain)
		case "hello":
			handleHello()
		case "exit":
			os.Exit(0)
		default:
			fmt.Println("Invalid command. Try again.")
		}
	}
}

func handleAddTransaction(data string, chain *blockchain.Blockchain) {
	// Thêm giao dịch và khai thác một khối mới
	chainMutex.Lock()
	defer chainMutex.Unlock()

	latestBlock := chain.GetLatestBlock()
	transaction := &blockchain.Transaction{Data: []byte(data)}
	latestBlock.Transactions = append(latestBlock.Transactions, transaction)

	newBlock := blockchain.GenerateBlock(latestBlock, latestBlock.Transactions)
	chain.AddBlock(newBlock.Transactions)

	fmt.Println("Transaction added and block mined successfully.")
}

func handlePrintChain(chain *blockchain.Blockchain) {
	// In ra thông tin chuỗi blockchain
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

func handleHello() {
	fmt.Println("Hello from node!")
}

func handleClient(conn net.Conn, chain *blockchain.Blockchain) {
	fmt.Printf("Client connected: %s\n", conn.RemoteAddr().String())

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		receivedMsg := scanner.Text()
		fmt.Printf("Received message from client: %s\n", receivedMsg)

		// Xử lý lệnh từ client (nếu cần)
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}
