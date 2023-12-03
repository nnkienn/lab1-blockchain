package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"github.com/nnkienn/lab1-blockchain/network"

)

var nodes = []string{"127.0.0.1:3001", "127.0.0.1:3002", "127.0.0.1:3003"}
var chainMutex sync.Mutex

func main() {
	chain := block.BlockChain{}

	for _, node := range nodes {
		go startServer(node, &chain)
	}

	go startClient(&chain)

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

func startClient(chain *block.BlockChain) {
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
			addTransaction(chain)
		case "printchain":
			printBlockchain(chain)
		case "hello":
			sendHelloRequest(clientPort)
		default:
			fmt.Println("Invalid command")
		}
	}
}

func addTransaction(chain *block.BlockChain) {
	fmt.Print("Enter transaction data: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	data := scanner.Text()

	transaction := &block.Transaction{
		Data: []byte(data),
	}

	chainMutex.Lock()
	latestBlock := chain.GetLatestBlock()
	latestBlock.AddTransaction(transaction.Data)
	chainMutex.Unlock()

	fmt.Println("Transaction added.")
}

func printBlockchain(chain *block.BlockChain) {
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

func handleClient(conn net.Conn, chain *block.BlockChain) {
	fmt.Printf("Client connected: %s\n", conn.RemoteAddr().String())

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		receivedMsg := scanner.Text()
		fmt.Printf("Received message from client: %s\n", receivedMsg)

		handleCommand(receivedMsg, chain, conn)
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}
