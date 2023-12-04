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

type Blockchain struct {
	Blocks []*Block
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

func CalculateMerkleRoot(transactions []*Transaction) []byte {
	var hashes [][]byte

	for _, transaction := range transactions {
		hashTransaction := sha256.Sum256(transaction.Data)
		hashes = append(hashes, hashTransaction[:])
	}

	for len(hashes) > 1 {
		var levelHashes [][]byte
		for i := 0; i < len(hashes)-1; i += 2 {
			combined := append(hashes[i], hashes[i+1]...)
			hash := sha256.Sum256(combined)
			levelHashes = append(levelHashes, hash[:])
		}
		if len(hashes)%2 != 0 {
			levelHashes = append(levelHashes, hashes[len(hashes)-1])
		}
		hashes = levelHashes
	}

	return hashes[0]
}

func (chain *Blockchain) AddBlock(transactions []*Transaction) {
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

func (chain *Blockchain) PrintBlockchain() {
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

func (chain *Blockchain) BuildMerkleTree() *MerkleTree {
	var transactions []*Transaction

	for _, block := range chain.Blocks {
		transactions = append(transactions, block.Transactions...)
	}

	merkleRoot := CalculateMerkleRoot(transactions)

	return NewMerkleTree(transactions, merkleRoot)
}

func (chain *Blockchain) CheckTransactionInMerkleTree(transactionData string) bool {
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
