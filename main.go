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

// P2PNode represents a node in the P2P model
type P2PNode struct {
	ServerAddress string
	BlockChain    *block.BlockChain
	Peers         map[string]struct{} // Keep track of connected peers
}

// StartP2PServer starts a P2P server for the node
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

// handleP2PConnection handles P2P connection to the node
func (node *P2PNode) handleP2PConnection(conn net.Conn) {
	defer conn.Close()

	peerAddress := conn.RemoteAddr().String()
	node.Peers[peerAddress] = struct{}{}

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		commandStr := scanner.Text()
		fmt.Println("Received P2P command from", peerAddress+":", commandStr)

		command, args := parseP2PCommand(commandStr)
		node.executeP2PCommand(command, args, conn)
	}
}

// parseP2PCommand parses a P2P command from a string
func parseP2PCommand(commandStr string) (string, []string) {
	parts := strings.Split(commandStr, "|")
	return parts[0], parts[1:]
}

// executeP2PCommand executes a P2P command
func (node *P2PNode) executeP2PCommand(command string, args []string, conn net.Conn) {
	mutex.Lock()
	defer mutex.Unlock()

	switch command {
	case "ADD_BLOCK":
		node.addBlockToChain(args)
	case "REQUEST_CHAIN":
		node.sendBlockChain(conn)
	case "RECEIVE_CHAIN":
		node.receiveBlockChainFromPeer(conn.RemoteAddr().String(), args)
	case "QUERY_TRANSACTION":
		node.queryTransaction(args[0], conn)
	default:
		fmt.Println("Unknown P2P command:", command)
	}
}

// addBlockToChain adds a new block to the node's blockchain
func (node *P2PNode) addBlockToChain(args []string) {
	if len(args) != 4 {
		fmt.Println("Invalid ADD_BLOCK command format")
		return
	}

	timestamp, prevBlockHash, merkleRoot, transactions := args[0], args[1], args[2], strings.Split(args[3], ",")
	timestampInt, _ := strconv.ParseInt(timestamp, 10, 64)

	newBlock := &block.Block{
		Timestamp:     timestampInt,
		PrevBlockHash: []byte(prevBlockHash),
		MerkleRoot:    []byte(merkleRoot),
		Transactions:  createTransactions(transactions),
	}

	node.BlockChain.AddBlock(newBlock.Transactions)
	node.sendBlockChainToPeers()
}

// createTransactions creates a list of transactions from input strings
func createTransactions(transactionStrings []string) []*block.Transaction {
	var transactions []*block.Transaction
	for _, transactionString := range transactionStrings {
		transactions = append(transactions, &block.Transaction{Data: []byte(transactionString)})
	}
	return transactions
}

// sendBlockChain sends the node's blockchain to a P2P connection
func (node *P2PNode) sendBlockChain(conn net.Conn) {
	chain := node.BlockChain
	chainString := serializeBlockchain(chain)
	conn.Write([]byte("RECEIVE_CHAIN|" + chainString + "\n"))
}

// serializeBlockchain converts blockchain to a string
func serializeBlockchain(chain *block.BlockChain) string {
	var result string
	for _, blk := range chain.Blocks {
		result += serializeBlock(blk) + "|"
	}
	return result
}

// serializeBlock converts a block to a string
func serializeBlock(blk *block.Block) string {
	return fmt.Sprintf("%d|%x|%x|%x|", blk.Timestamp, blk.PrevBlockHash, blk.MerkleRoot, serializeTransactions(blk.Transactions))
}

// serializeTransactions converts a list of transactions to a string
// serializeTransactions chuyển đổi một danh sách giao dịch thành một chuỗi
func serializeTransactions(transactions []*block.Transaction) string {
	var result string
	for _, transaction := range transactions {
		result += string(transaction.Data) + ","
	}
	if len(result) > 0 {
		result = result[:len(result)-1] // Loại bỏ dấu phẩy ở cuối
	}
	return result
}


// receiveBlockChainFromPeer receives a blockchain from a peer and updates the node's blockchain
// receiveBlockChainFromPeer receives a blockchain from a peer and updates the node's blockchain
func (node *P2PNode) receiveBlockChainFromPeer(peerAddress string, blocks []string) {
	for _, blockString := range blocks {
		timestamp, prevBlockHash, merkleRoot, transactions := parseBlockString(blockString)

		// Convert timestamp to string
		timestampStr := strconv.FormatInt(timestamp, 10)

		// Pass data as a single slice to addBlockToChain
		node.addBlockToChain([]string{timestampStr, prevBlockHash, merkleRoot, strings.Join(transactions, ",")})
	}

	fmt.Printf("Blockchain updated from peer %s\n", peerAddress)
}

// parseBlockString phân tích một chuỗi thành các giá trị cần thiết cho một block
func parseBlockString(blockString string) (int64, string, string, []string) {
	parts := strings.SplitN(blockString, "|", 4)
	if len(parts) != 4 {
		// Xử lý trường hợp chuỗi khối không chứa đủ phần
		fmt.Println("Định dạng chuỗi khối không hợp lệ:", blockString)
		return 0, "", "", nil
	}
	timestamp, prevBlockHash, merkleRoot := parts[0], parts[1], parts[2]
	transactions := strings.Split(parts[3], ",")
	return parseTimestamp(timestamp), prevBlockHash, merkleRoot, transactions
}


// parseTimestamp parses a timestamp string into an int64 value
func parseTimestamp(timestamp string) int64 {
	timestampInt, _ := strconv.ParseInt(timestamp, 10, 64)
	return timestampInt
}

// queryTransaction queries a transaction based on the Merkle tree
func (node *P2PNode) queryTransaction(transactionData string, conn net.Conn) {
	var response string
	if node.BlockChain != nil {
		exists := node.BlockChain.CheckTransactionInMerkleTree(transactionData)
		response = "TRANSACTION_QUERY_RESULT|" + strconv.FormatBool(exists)
	} else {
		response = "TRANSACTION_QUERY_RESULT|false"
	}

	if conn != nil {
		conn.Write([]byte(response + "\n"))
	} else {
		fmt.Println(response)
	}
}


// main
func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <port>")
		return
	}

	port := os.Args[1]
	node := &P2PNode{
		ServerAddress: "localhost:" + port,
		BlockChain:    &block.BlockChain{Blocks: []*block.Block{}},
		Peers:         make(map[string]struct{}),
	}

	go node.StartP2PServer()

	fmt.Println("P2P server is running on port", port+". Enter commands:")

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
			node.sendBlockChainToPeers()
		case "MINE_BLOCK":
			node.mineBlock()
			node.sendBlockChainToPeers()
		case "CONNECT":
			fmt.Print("Enter address to connect (e.g., localhost:3001): ")
			address, _ := reader.ReadString('\n')
			address = strings.TrimSpace(address)
			go node.connectToPeer(address)
		case "QUERY_TRANSACTION":
			fmt.Print("Enter transaction data to query: ")
			data, _ := reader.ReadString('\n')
			data = strings.TrimSpace(data)
			node.queryTransaction(data, nil)
		default:
			fmt.Println("Unknown command:", command)
		}
	}
}
// addTransactionToMempool adds a new transaction to the node's mempool
func (node *P2PNode) addTransactionToMempool(data string) {
	transaction := &block.Transaction{Data: []byte(data)}
	node.BlockChain.AddTransactionToMempool(transaction)
	fmt.Println("Transaction added to mempool.")
}

// mineBlock mines a block from the mempool and adds it to the blockchain
func (node *P2PNode) mineBlock() {
	node.BlockChain.MineBlock()
	fmt.Println("Block mined and added to the blockchain.")
	node.sendBlockChainToPeers()
}

// connectToPeer connects to a peer address
func (node *P2PNode) connectToPeer(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error connecting to peer:", err)
		return
	}
	defer conn.Close()

	// Add the connected peer to the list
	node.Peers[address] = struct{}{}

	fmt.Printf("Connected to peer at %s\n", address)

	// Request the peer's blockchain
	conn.Write([]byte("REQUEST_CHAIN\n"))
}

// sendBlockChainToPeers sends the node's blockchain to all connected peers
func (node *P2PNode) sendBlockChainToPeers() {
	for peer := range node.Peers {
		node.sendBlockChainToPeer(peer)
	}
}

// sendBlockChainToPeer sends the node's blockchain to a specific peer
func (node *P2PNode) sendBlockChainToPeer(peerAddress string) {
	conn, err := net.Dial("tcp", peerAddress)
	if err != nil {
		fmt.Println("Error connecting to peer:", err)
		return
	}
	defer conn.Close()

	chain := node.BlockChain
	chainString := serializeBlockchain(chain)
	conn.Write([]byte("RECEIVE_CHAIN|" + chainString + "\n"))
}

// printBlockchain prints information about the blockchain
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
