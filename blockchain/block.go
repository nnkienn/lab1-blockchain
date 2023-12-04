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
	MerkleRoot    []byte // Thêm trường MerkleRoot
}

type Transaction struct {
	Data []byte
}

type BlockChain struct {
	Blocks []*Block
}

// MerkleNode đại diện cho một node trong cây Merkle
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

// MerkleTree đại diện cho cây Merkle
type MerkleTree struct {
	Root *MerkleNode
}

// NewMerkleTree tạo và trả về một cây Merkle mới
func NewMerkleTree(transactions []*Transaction, merkleRoot []byte) *MerkleTree {
	var nodes []*MerkleNode

	// Tạo các nút lá
	for _, transaction := range transactions {
		nodes = append(nodes, &MerkleNode{Data: transaction.Data})
	}

	// Xây dựng cây Merkle từ nút lá
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
			// Nếu số lượng nút là số lẻ, thì thêm nút cuối cùng
			newLevel = append(newLevel, nodes[len(nodes)-1])
		}

		nodes = newLevel
	}

	return &MerkleTree{Root: nodes[0]}
}

func (block *Block) setHash() {
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

	newBlock.MerkleRoot = CalculateMerkleRoot(transactions) // Tính toán MerkleRoot
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

func (chain *BlockChain) BuildMerkleTree() *MerkleTree {
	var transactions []*Transaction

	for _, block := range chain.Blocks {
		transactions = append(transactions, block.Transactions...)
	}

	merkleRoot := CalculateMerkleRoot(transactions)

	// Tạo và trả về cây Merkle
	return NewMerkleTree(transactions, merkleRoot)
}

// CheckTransactionInMerkleTree kiểm tra xem một giao dịch có tồn tại trong cây Merkle hay không
// CheckTransactionInMerkleTree kiểm tra xem một giao dịch có tồn tại trong cây Merkle của blockchain hay không
func (chain *BlockChain) CheckTransactionInMerkleTree(transactionData string) bool {
    var transactions []*Transaction

    // Tập hợp tất cả các giao dịch trong blockchain
    for _, block := range chain.Blocks {
        transactions = append(transactions, block.Transactions...)
    }

    // Tính toán Merkle root cho tất cả các giao dịch
    merkleRoot := CalculateMerkleRoot(transactions)

    // Tạo cây Merkle cho tất cả các giao dịch
    merkleTree := NewMerkleTree(transactions, merkleRoot)

    // Kiểm tra xem giao dịch có tồn tại trong cây Merkle hay không
    return merkleTree.CheckTransaction(&Transaction{Data: []byte(transactionData)})
}


// CheckTransaction kiểm tra xem một giao dịch có tồn tại trong cây Merkle hay không
func (tree *MerkleTree) CheckTransaction(transaction *Transaction) bool {
	return tree.checkTransaction(tree.Root, transaction.Data)
}

// checkTransaction kiểm tra xem một giao dịch có tồn tại trong cây Merkle hay không (hàm đệ quy)
func (tree *MerkleTree) checkTransaction(node *MerkleNode, transactionData []byte) bool {
	if node == nil {
		return false
	}

	// Nếu là nút lá, so sánh dữ liệu giao dịch
	if node.Left == nil && node.Right == nil {
		return bytes.Equal(node.Data, transactionData)
	}

	// Nếu không phải nút lá, kiểm tra ở cả hai phía của cây
	leftResult := tree.checkTransaction(node.Left, transactionData)
	rightResult := tree.checkTransaction(node.Right, transactionData)

	return leftResult || rightResult
}
