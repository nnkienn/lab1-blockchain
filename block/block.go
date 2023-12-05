// block.go

package block

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"
)

type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	MerkleRoot    []byte
}

type Transaction struct {
	Data []byte
}

// BlockChain represents a blockchain
type BlockChain struct {
	Blocks  []*Block
	Mempool []*Transaction
}

// GetMempoolTransactions returns transactions from the mempool
func (chain *BlockChain) GetMempoolTransactions() []*Transaction {
	return chain.Mempool
}

// AddTransactionToMempool adds a transaction to the mempool
func (chain *BlockChain) AddTransactionToMempool(transaction *Transaction) {
	chain.Mempool = append(chain.Mempool, transaction)
}

// MineBlock mines a block with transactions from the mempool
func (chain *BlockChain) MineBlock() {
	blockData := chain.GetMempoolTransactions()
	chain.AddBlock(blockData)
	chain.Mempool = []*Transaction{} // Clear the mempool after mining
}

// MerkleNode represents a node in the Merkle tree
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

// MerkleTree represents a Merkle tree
type MerkleTree struct {
	Root *MerkleNode
}

// NewMerkleTree creates and returns a new Merkle tree
func NewMerkleTree(transactions []*Transaction, merkleRoot []byte) *MerkleTree {
	var nodes []*MerkleNode

	for _, transaction := range transactions {
		nodes = append(nodes, &MerkleNode{Data: transaction.Data})
	}

	for len(nodes) > 1 {
		var newLevel []*MerkleNode

		for i := 0; i < len(nodes)-1; i += 2 {
			combinedData := append(nodes[i].Data, nodes[i+1].Data...)
			hash := sha256.Sum256(combinedData)

			newNode := &MerkleNode{
				Left:  nodes[i],
				Right: nodes[i+1],
				Data:  hash[:],
			}

			newLevel = append(newLevel, newNode)
		}

		if len(nodes)%2 != 0 {
			newLevel = append(newLevel, nodes[len(nodes)-1])
		}

		nodes = newLevel
	}

	return &MerkleTree{Root: nodes[0]}
}

func (block *Block) SetHash() {
	headers := []byte(fmt.Sprintf("%x%x%d", block.PrevBlockHash, block.MerkleRoot, block.Timestamp))
	hash := sha256.Sum256(headers)
	block.Hash = hash[:]
}

// CalculateMerkleRoot tính toán Merkle Root từ danh sách giao dịch
func CalculateMerkleRoot(transactions []*Transaction) []byte {
	if len(transactions) == 0 {
		return nil
	}

	var merkleTree [][]byte

	for _, tx := range transactions {
		merkleTree = append(merkleTree, tx.Data)
	}

	for len(merkleTree) > 1 {
		var newMerkleTree [][]byte
		for i := 0; i < len(merkleTree); i += 2 {
			left := merkleTree[i]
			right := merkleTree[i+1]
			concatenation := append(left, right...)
			hash := sha256.Sum256(concatenation)
			newMerkleTree = append(newMerkleTree, hash[:])
		}
		// Nếu số lượng node là lẻ, duplicate node cuối cùng
		if len(merkleTree)%2 != 0 {
			lastNode := merkleTree[len(merkleTree)-1]
			hash := sha256.Sum256(append(lastNode, lastNode...))
			newMerkleTree = append(newMerkleTree, hash[:])
		}
		merkleTree = newMerkleTree
	}

	return merkleTree[0]
}

func (chain *BlockChain) AddBlock(transactions []*Transaction) {
	var preBlockHash []byte

	chainSize := len(chain.Blocks)
	if chainSize > 0 {
		preBlockHash = chain.Blocks[chainSize-1].Hash
	}

	newBlock := &Block{
		Timestamp:     time.Now().Unix(),
		PrevBlockHash: preBlockHash,
		Transactions:  transactions,
	}

	newBlock.MerkleRoot = CalculateMerkleRoot(transactions)
	newBlock.SetHash()

	chain.Blocks = append(chain.Blocks, newBlock)
}

func (chain *BlockChain) PrintBlockchain() {
	for _, block := range chain.Blocks {
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Merkle root: %x\n", block.MerkleRoot)
		fmt.Println("Transactions:")
		for _, transaction := range block.Transactions {
			fmt.Printf("- %s\n", string(transaction.Data))
		}
		fmt.Printf("%x\n", block.Hash)
		fmt.Println("--------------")
	}
}

func (chain *BlockChain) BuildMerkleTree() *MerkleTree {
	var transactions []*Transaction

	for _, block := range chain.Blocks {
		transactions = append(transactions, block.Transactions...)
	}

	merkleRoot := CalculateMerkleRoot(transactions)

	return NewMerkleTree(transactions, merkleRoot)
}

func (chain *BlockChain) CheckTransactionInMerkleTree(transactionData string) bool {
	var transactions []*Transaction

	for _, block := range chain.Blocks {
		transactions = append(transactions, block.Transactions...)
	}

	merkleRoot := CalculateMerkleRoot(transactions)

	merkleTree := NewMerkleTree(transactions, merkleRoot)

	return merkleTree.CheckTransaction(&Transaction{Data: []byte(transactionData)})
}

func (tree *MerkleTree) CheckTransaction(transaction *Transaction) bool {
	return tree.checkTransaction(tree.Root, transaction.Data)
}

func (tree *MerkleTree) checkTransaction(node *MerkleNode, transactionData []byte) bool {
	if node == nil {
		return false
	}

	if node.Left == nil && node.Right == nil {
		return bytes.Equal(node.Data, transactionData)
	}

	leftResult := tree.checkTransaction(node.Left, transactionData)
	rightResult := tree.checkTransaction(node.Right, transactionData)

	return leftResult || rightResult
}
