package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

func receiveResponse(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		response := scanner.Text()
		fmt.Println("Server response:", response)

		// Tách Merkle root từ dòng phản hồi
		if strings.HasPrefix(response, "Merkle root:") {
			parts := strings.Split(response, ":")
			serverMerkleRoot := strings.TrimSpace(parts[1])
			if serverMerkleRoot == "Transaction does not exist in the Merkle tree." {
				fmt.Println("Transaction does not exist in the Merkle tree.")
			} else {
				fmt.Println("Transaction exists in the Merkle tree.")
			}
		}
	}
}

func sendCommand(conn net.Conn, command string) {
	_, err := conn.Write([]byte(command + "\n"))
	if err != nil {
		fmt.Println("Error sending command:", err)
	}
}

func createBlock(conn net.Conn, transactions []string) {
	// Gửi lệnh ADD_BLOCK|transaction1,transaction2,...
	command := fmt.Sprintf("ADD_BLOCK|%s", strings.Join(transactions, ","))
	sendCommand(conn, command)
}

func queryTransaction(conn net.Conn, transactionData string) {
	// Gửi yêu cầu để tạo Merkle Tree từ blockchain
	sendCommand(conn, "BUILD_MERKLE_TREE")

	// Nhận Merkle Root từ server
	scanner := bufio.NewScanner(conn)
	var merkleRoot string
	for scanner.Scan() {
		response := scanner.Text()
		fmt.Println("Server response:", response)

		// Tách Merkle root từ dòng phản hồi
		if strings.HasPrefix(response, "Merkle root:") {
			parts := strings.Split(response, ":")
			merkleRoot = strings.TrimSpace(parts[1])
			break
		}
	}

	// Gửi lệnh để kiểm tra giao dịch trong Merkle Tree
	command := fmt.Sprintf("QUERY_TRANSACTION|%s|%s", merkleRoot, transactionData)
	sendCommand(conn, command)
}

func main() {
	// Sử dụng flag để nhận giá trị cổng từ dòng lệnh
	port := flag.Int("port", 8080, "Server port")
	flag.Parse()

	serverAddress := fmt.Sprintf("localhost:%d", *port)

	// Kết nối đến máy chủ
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Printf("Connected to server on port %d.\n", *port)

	go receiveResponse(conn)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter command: ")
		scanner.Scan()
		command := scanner.Text()

		if strings.HasPrefix(command, "CREATE_BLOCK") {
			// Example: CREATE_BLOCK|transaction1,transaction2,transaction3
			parts := strings.Split(command, "|")
			transactions := strings.Split(parts[1], ",")
			createBlock(conn, transactions)
		} else if strings.HasPrefix(command, "QUERY_TRANSACTION") {
			// Example: QUERY_TRANSACTION|transaction_data
			parts := strings.Split(command, "|")
			transactionData := parts[1]
			queryTransaction(conn, transactionData)
		} else {
			sendCommand(conn, command)
		}
	}

	fmt.Println("Connection closed.")
}
