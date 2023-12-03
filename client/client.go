package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/nnkienn/lab1-blockchain/server/blockchain"
)

var nodes = []string{"127.0.0.1:3001", "127.0.0.1:3002", "127.0.0.1:3003"}
var chainMutex sync.Mutex

func main() {
	chain := blockchain.NewBlockchain()

	for _, node := range nodes {
		go startServer(node, chain)
	}

	select {}
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

func handleCommand(command string, chain *blockchain.Blockchain, conn net.Conn) {
	parts := strings.Fields(command)
	switch parts[0] {
	case "printchain":
		sendBlockchainInfo(chain, conn)
	case "hello":
		sendHelloResponse(conn)
	case "addtransaction":
		if len(parts) != 2 {
			fmt.Fprintln(conn, "Invalid addtransaction command. Usage: addtransaction <transaction_data>")
			return
		}
		data := parts[1]
		addTransactionAndMineBlock(chain, conn, data)
	case "addnode":
		addNode(parts, conn)
	default:
		fmt.Fprintln(conn, "Invalid command")
	}
}

func addTransactionAndMineBlock(chain *blockchain.Blockchain, conn net.Conn, data string) {
	chainMutex.Lock()
	defer chainMutex.Unlock()

	latestBlock := chain.GetLatestBlock()

	// Tạo một giao dịch mới và thêm vào khối hiện tại
	transaction := &blockchain.Transaction{Data: []byte(data)}
	latestBlock.Transactions = append(latestBlock.Transactions, transaction)

	// Đối với mục đích minh họa, bạn có thể thêm điều kiện để chỉ khai thác khi đạt đến số lượng giao dịch mong muốn.
	// Ví dụ: if len(latestBlock.Transactions) >= 2 { ... }
	// Trong thực tế, điều này sẽ được xử lý thông qua các chính sách khai thác thích hợp.

	// Tạo khối mới và thêm vào chuỗi blockchain
	transactions := append([]*blockchain.Transaction{}, latestBlock.Transactions...)
	newBlock := blockchain.GenerateBlock(latestBlock, transactions)
	chain.AddBlock(newBlock.Transactions)

	// Gửi thông báo về việc thêm giao dịch và khối thành công
	fmt.Fprintln(conn, "Transaction added and block mined successfully.")
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

func handleClient(conn net.Conn, chain *blockchain.Blockchain) {
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

func sendBlockchainInfo(chain *blockchain.Blockchain, conn net.Conn) {
	chainMutex.Lock()
	defer chainMutex.Unlock()

	for _, block := range chain.GetBlocks() {
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
