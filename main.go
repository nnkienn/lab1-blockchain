// main.go

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

// handleP2PConnection handles a P2P connection to the node
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

// parseP2PCommand parses a P2P command from a string
func parseP2PCommand(commandStr string) (string, []string) {
	parts := strings.Split(commandStr, "|")
	command := parts[0]
	args := parts[1:]
	return command, args
}

// executeP2PCommand executes a P2P command
func (node *P2PNode) executeP2PCommand(command string, args []string, conn net.Conn) {
	switch command {
	case "ADD_BLOCK":
		timestamp := args[0]
		prevBlockHash := args[1]
		merkleRoot := args[2]
		transactions := strings.Split(args[3], ",")

		node.addBlockToChain(timestamp, prevBlockHash, merkleRoot, transactions)
	case "REQUEST_CHAIN":
		node.sendBlockChain(conn)
	case "RECEIVE_CHAIN":
		blocks := args
		node.receiveBlockChain(blocks)
	default:
		fmt.Println("Unknown P2P command:", command)
	}
}

// addBlockToChain adds a new block to the node's blockchain
func (node *P2PNode) addBlockToChain(timestamp, prevBlockHash, merkleRoot string, transactions []string) {
	timestampInt, _ := strconv.ParseInt(timestamp, 10, 64)

	newBlock := &block.Block{
		Timestamp:     timestampInt,
		PrevBlockHash: []byte(prevBlockHash),
		MerkleRoot:    []byte(merkleRoot),
		Transactions:  createTransactions(transactions),
	}

	transactionData := []string{}
	for _, transaction := range newBlock.Transactions {
		transactionData = append(transactionData, string(transaction.Data))
	}

	node.BlockChain.AddBlock(createTransactions(transactionData))
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
func serializeTransactions(transactions []*block.Transaction) string {
	var result string
	for _, transaction := range transactions {
		result += string(transaction.Data) + ","
	}
	return result[:len(result)-1]
}

// receiveBlockChain receives a blockchain from a string and updates the node's blockchain
func (node *P2PNode) receiveBlockChain(blocks []string) {
	for _, blockString := range blocks {
		timestamp, prevBlockHash, merkleRoot, transactions := parseBlockString(blockString)
		timestampStr := strconv.FormatInt(timestamp, 10)
		node.addBlockToChain(timestampStr, prevBlockHash, merkleRoot, transactions)
	}
}

// parseBlockString parses a string into values needed for a block
func parseBlockString(blockString string) (int64, string, string, []string) {
	parts := strings.Split(blockString, "|")
	timestamp, prevBlockHash, merkleRoot := parts[0], parts[1], parts[2]
	transactions := strings.Split(parts[3], ",")
	return parseTimestamp(timestamp), prevBlockHash, merkleRoot, transactions
}

// parseTimestamp parses a timestamp string into an int64 value
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

	go node.StartP2PServer()

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

// addTransactionToMempool adds a new transaction to the node's mempool
func (node *P2PNode) addTransactionToMempool(data string) {
	transaction := &block.Transaction{Data: []byte(data)}
	node.BlockChain.AddTransactionToMempool(transaction)
	fmt.Println("Transaction added to mempool.")
}

// mineBlock mines a block from the mempool and adds it to the blockchain
func (node *P2PNode) mineBlock() {
	blockData := node.BlockChain.GetMempoolTransactions()
	node.BlockChain.MineBlock(blockData)
	fmt.Println("Block mined and added to the blockchain.")
}

// connectToPeer connects to a peer address
func (node *P2PNode) connectToPeer(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error connecting to peer:", err)
		return
	}
	defer conn.Close()

	fmt.Printf("Connected to peer at %s\n", address)

	conn.Write([]byte("REQUEST_CHAIN\n"))
}

// printBlockchain prints the information of the blockchain
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
