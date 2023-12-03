package blockchain

import (
	"crypto/sha256"
	"time"
	"fmt"
)

type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	MerkleRoot    []byte // Thêm MerkleRoot vào cấu trúc Block
}

type Transaction struct {
	Data []byte
}

type BlockChain struct {
	Blocks []*Block
}

func (block *Block) setHash() {
	headers := []byte(string(block.PrevBlockHash) + string(block.MerkleRoot) + string(block.Timestamp)) // Sử dụng MerkleRoot thay vì HashTransactions
	hash := sha256.Sum256(headers)
	block.Hash = hash[:]
}

func HashTransactions(transactions []*Transaction) []byte {
	var hashTransactions []byte

	for _, transaction := range transactions {
		hashTransaction := sha256.Sum256(transaction.Data)
		hashTransactions = append(hashTransactions, hashTransaction[:]...)
	}

	finalHash := sha256.Sum256(hashTransactions)

	return finalHash[:]
}

func calculateMerkleRoot(transactions []*Transaction) []byte {
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

	newBlock.MerkleRoot = calculateMerkleRoot(transactions) // Tính toán MerkleRoot
	newBlock.setHash()

	chain.Blocks = append(chain.Blocks, newBlock)
}

func PrintBlockchain(chain *BlockChain) {
	// In ra thông tin của blockchain
	for _, block := range chain.Blocks {
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Merkle root: %x\n", block.MerkleRoot) // In ra MerkleRoot
		fmt.Println("Transactions:")
		for _, transaction := range block.Transactions {
			fmt.Printf("- %s\n", string(transaction.Data))
		}
		fmt.Printf("%x\n", block.Hash)
	}
}
