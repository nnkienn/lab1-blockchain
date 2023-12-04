package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/nnkienn/lab1-blockchain/block"
)

var mutex = &sync.Mutex{}

// P2PNode đại diện cho một nút trong mô hình P2P
	type P2PNode struct {
		ServerAddress string
		BlockChain    *block.BlockChain
	}

// StartP2PServer khởi động một server P2P cho nút
func (node *P2PNode) StartP2PServer() {
	listener, err := net.Listen("tcp", node.ServerAddress)
	if err != nil {
		fmt.Println("Error starting P2P server:", err)
		return
	}
	defer listener.Close()

	fmt.Printf("P2P server started on %s\n", node.ServerAddress)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go node.handleP2PConnection(conn)
	}
}

// handleP2PConnection xử lý kết nối P2P đến nút
func (node *P2PNode) handleP2PConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		commandStr := scanner.Text()
		fmt.Println("Received P2P command:", commandStr)

		command, args := parseP2PCommand(commandStr)
		node.executeP2PCommand(command, args, conn)
	}
}

// parseP2PCommand phân tích một lệnh P2P từ một chuỗi
func parseP2PCommand(commandStr string) (string, []string) {
	parts := strings.Split(commandStr, "|")
	command := parts[0]
	args := parts[1:]
	return command, args
}

// executeP2PCommand thực thi một lệnh P2P
func (node *P2PNode) executeP2PCommand(command string, args []string, conn net.Conn) {
	switch command {
	case "ADD_BLOCK":
		// Example: ADD_BLOCK|timestamp|prevBlockHash|merkleRoot|transaction1,transaction2,transaction3
		timestamp := args[0]
		prevBlockHash := args[1]
		merkleRoot := args[2]
		transactions := strings.Split(args[3], ",")

		node.addBlockToChain(timestamp, prevBlockHash, merkleRoot, transactions)
	case "REQUEST_CHAIN":
		// Example: REQUEST_CHAIN
		node.sendBlockChain(conn)
	case "RECEIVE_CHAIN":
		// Example: RECEIVE_CHAIN|block1|block2|...
		blocks := args
		node.receiveBlockChain(blocks)
	default:
		fmt.Println("Unknown P2P command:", command)
	}
}

// addBlockToChain thêm một block mới vào blockchain của nút
// addBlockToChain thêm một block mới vào blockchain của nút
func (node *P2PNode) addBlockToChain(timestamp, prevBlockHash, merkleRoot string, transactions []string) {
	// Convert string parameters to appropriate types
	timestampInt, _ := strconv.ParseInt(timestamp, 10, 64)

	// Create Block and add to the blockchain
	newBlock := &block.Block{
		Timestamp:     timestampInt,
		PrevBlockHash: []byte(prevBlockHash),
		MerkleRoot:    []byte(merkleRoot),
		Transactions:  createTransactions(transactions),
	}

	// Add the block to the blockchain
	node.BlockChain.AddBlock([]*block.Block{newBlock})
}

// createTransactions tạo một danh sách các giao dịch từ các chuỗi đầu vào
func createTransactions(transactionStrings []string) []*block.Transaction {
	var transactions []*block.Transaction
	for _, transactionString := range transactionStrings {
		transactions = append(transactions, &block.Transaction{Data: []byte(transactionString)})
	}
	return transactions
}

// sendBlockChain gửi blockchain của nút đến một kết nối P2P
func (node *P2PNode) sendBlockChain(conn net.Conn) {
	chain := node.BlockChain
	chainString := serializeBlockchain(chain)
	conn.Write([]byte("RECEIVE_CHAIN|" + chainString + "\n"))
}

// serializeBlockchain chuyển đổi blockchain thành một chuỗi
func serializeBlockchain(chain *block.BlockChain) string {
	var result string
	for _, blk := range chain.Blocks {
		result += serializeBlock(blk) + "|"
	}
	return result
}

// serializeBlock chuyển đổi một block thành một chuỗi
func serializeBlock(blk *block.Block) string {
	return fmt.Sprintf("%d|%x|%x|%x|", blk.Timestamp, blk.PrevBlockHash, blk.MerkleRoot, serializeTransactions(blk.Transactions))
}

// serializeTransactions chuyển đổi danh sách giao dịch thành một chuỗi
func serializeTransactions(transactions []*block.Transaction) string {
	var result string
	for _, transaction := range transactions {
		result += string(transaction.Data) + ","
	}
	return result[:len(result)-1] // Remove the trailing comma
}

// receiveBlockChain nhận blockchain từ một chuỗi và cập nhật blockchain của nút
func (node *P2PNode) receiveBlockChain(blocks []string) {
	for _, blockString := range blocks {
		timestamp, prevBlockHash, merkleRoot, transactions := parseBlockString(blockString)

		// Convert timestamp to string
		timestampStr := strconv.FormatInt(timestamp, 10)

		node.addBlockToChain(timestampStr, prevBlockHash, merkleRoot, transactions)
	}
}

// parseBlockString phân tích một chuỗi thành các giá trị cần thiết cho một block
func parseBlockString(blockString string) (int64, string, string, []string) {
	parts := strings.Split(blockString, "|")
	timestamp, prevBlockHash, merkleRoot := parts[0], parts[1], parts[2]
	transactions := strings.Split(parts[3], ",")
	return parseTimestamp(timestamp), prevBlockHash, merkleRoot, transactions
}

// parseTimestamp phân tích một chuỗi timestamp thành một giá trị int64
func parseTimestamp(timestamp string) int64 {
	timestampInt, _ := strconv.ParseInt(timestamp, 10, 64)
	return timestampInt
}

// main
func main() {
	node := &P2PNode{
		ServerAddress: "localhost:3000",
		BlockChain:    &block.BlockChain{Blocks: []*block.Block{}},
	}

	// Start P2P server
	go node.StartP2PServer()

	// Interact with the user through the console
	fmt.Println("P2P server is running. Enter commands:")

	for {
		reader := bufio.NewReader(os.Stdin)
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		switch command {
		case "PRINT_CHAIN":
			printBlockchain(node.BlockChain)
		case "ADD_TRANSACTION":
			fmt.Print("Enter transaction data: ")
			data, _ := reader.ReadString('\n')
			data = strings.TrimSpace(data)
			node.addTransactionToMempool(data)
		case "MINE_BLOCK":
			node.mineBlock()
		case "CONNECT":
			fmt.Print("Enter address to connect (e.g., localhost:3001): ")
			address, _ := reader.ReadString('\n')
			address = strings.TrimSpace(address)
			go node.connectToPeer(address)
		default:
			fmt.Println("Unknown command:", command)
		}
	}
}

// addTransactionToMempool thêm một giao dịch mới vào mempool của nút
func (node *P2PNode) addTransactionToMempool(data string) {
	transaction := &block.Transaction{Data: []byte(data)}
	node.BlockChain.AddTransactionToMempool(transaction)
	fmt.Println("Transaction added to mempool.")
}

// mineBlock khai thác một block từ mempool và thêm vào blockchain
func (node *P2PNode) mineBlock() {
	blockData := node.BlockChain.GetMempoolTransactions()
	node.BlockChain.MineBlock(blockData)
	fmt.Println("Block mined and added to the blockchain.")
}

// connectToPeer kết nối đến một địa chỉ peer
func (node *P2PNode) connectToPeer(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error connecting to peer:", err)
		return
	}
	defer conn.Close()

	fmt.Printf("Connected to peer at %s\n", address)

	// Send a request for the peer's blockchain
	conn.Write([]byte("REQUEST_CHAIN\n"))
}

// printBlockchain in ra thông tin của blockchain
func printBlockchain(chain *block.BlockChain) {
	for _, blk := range chain.Blocks {
		fmt.Printf("Timestamp: %d\n", blk.Timestamp)
		fmt.Printf("Prev. hash: %x\n", blk.PrevBlockHash)
		fmt.Printf("Merkle root: %x\n", blk.MerkleRoot)
		fmt.Println("Transactions:")
		for _, transaction := range blk.Transactions {
			fmt.Printf("- %s\n", string(transaction.Data))
		}
		fmt.Printf("%x\n", blk.Hash)
		fmt.Println("--------------")
	}
}
