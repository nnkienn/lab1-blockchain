package block

import (
	"crypto/sha256"
	"time"
	"fmt"
)

type Block struct {
	Timestamp int64
	Transactions []*Transaction
	PrevBlockHash []byte
	Hash []byte
}

type Transaction struct {
	Data []byte
}

type BlockChain struct {
	Blocks []*Block
}

func (block *Block) setHash() {
	headers := []byte(string(block.PrevBlockHash) + string(HashTransactions(block.Transactions)) + string(block.Timestamp))
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

func (chain *BlockChain) AddBlock(transactions []*Transaction) {
	var preBlockHash []byte

	chain_size := len(chain.Blocks)
	if (chain_size > 0 ) {
		preBlockHash = chain.Blocks[chain_size - 1].Hash
	}

	newBlock := &Block{
		Timestamp: time.Now().Unix(),
		PrevBlockHash: preBlockHash,
		Transactions: transactions,
	}

	newBlock.setHash()

	chain.Blocks = append(chain.Blocks, newBlock)
}

func PrintBlockchain(chain *BlockChain) {
	// In ra thông tin của blockchain
	for _, block := range chain.Blocks {
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)

		fmt.Println("Transactions:")
		for _, transaction := range block.Transactions {
			fmt.Printf("- %s\n", string(transaction.Data))
		}

		fmt.Printf("%x\n", block.Hash)
	}
}
