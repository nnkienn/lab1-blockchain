package network

import (
	"github.com/nnkienn/lab1-blockchain/server/blockchain"
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"strings"
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

func handleCommand(command string, chain *block.BlockChain, conn net.Conn) {
	parts := strings.Fields(command)
	switch parts[0] {
	case "printchain":
		sendBlockchainInfo(chain, conn)
	case "hello":
		sendHelloResponse(conn)
	case "addtransaction":
		addTransactionAndMineBlock(chain, conn)
	case "addnode":
		addNode(parts, conn)
	default:
		fmt.Fprintln(conn, "Invalid command")
	}
}

func addNode(parts []string, conn net.Conn) {
	if len(parts) != 2 {
		fmt.Fprintln(conn, "Invalid addnode command. Usage: addnode <node_address>")
		return
	}

	newNode := parts[1]
	nodes = append(nodes, newNode)
	fmt.Fprintf(conn, "Node %s added to the network.\n", newNode)
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

func sendBlockchainInfo(chain *block.BlockChain, conn net.Conn) {
	chainMutex.Lock()
	defer chainMutex.Unlock()

	for _, block := range chain.Blocks {
		blockInfo := fmt.Sprintf("Timestamp: %d\nPrev. hash: %x\n", block.Timestamp, block.PrevBlockHash)
		for _, transaction := range block.Transactions {
			blockInfo += fmt.Sprintf("Transaction: %s\n", string(transaction.Data))
		}
		blockInfo += fmt.Sprintf("Hash: %x\n", block.Hash)

		fmt.Fprintln(conn, blockInfo)
	}
}

func sendHelloResponse(conn net.Conn) {
	fmt.Fprintln(conn, "Hello from node!")
}
