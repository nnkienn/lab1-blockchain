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

func queryTransaction(conn net.Conn, merkleRoot string) {
	// Gửi lệnh PRINT_BLOCKCHAIN để tra cứu Merkle tree
	sendCommand(conn, "PRINT_BLOCKCHAIN")

	// Chờ nhận phản hồi từ server
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		response := scanner.Text()
		fmt.Println(response)

		// Tách Merkle root từ dòng phản hồi
		if strings.HasPrefix(response, "Merkle root:") {
			parts := strings.Split(response, ":")
			serverMerkleRoot := strings.TrimSpace(parts[1])
			if serverMerkleRoot == merkleRoot {
				fmt.Println("Transaction exists in the Merkle tree.")
			} else {
				fmt.Println("Transaction does not exist in the Merkle tree.")
			}
			break
		}
	}
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
			// Example: QUERY_TRANSACTION|merkleRoot
			parts := strings.Split(command, "|")
			merkleRoot := parts[1]
			queryTransaction(conn, merkleRoot)
		} else {
			sendCommand(conn, command)
		}
	}

	fmt.Println("Connection closed.")
}
