package block

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// Transaction represents a transaction in the blockchain.
type Transaction struct {
	Data []byte
}

// Block represents a block in the blockchain.
type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
}

// NewBlock creates a new block.
func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		PrevBlockHash: prevBlockHash,
		Transactions:  transactions,
	}

	block.setHash()
	return block
}

// setHash calculates and sets the hash of the block.
func (block *Block) setHash() {
	headers := []byte(fmt.Sprintf("%d%s%s", block.Timestamp, block.PrevBlockHash, block.getTransactionsString()))
	hash := sha256.Sum256(headers)
	block.Hash = hash[:]
}

// getTransactionsString returns a concatenated string of transaction data.
func (block *Block) getTransactionsString() string {
	var transactionsData string
	for _, tx := range block.Transactions {
		transactionsData += string(tx.Data)
	}
	return transactionsData
}
