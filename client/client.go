package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

// client.go

// ...

func receiveResponse(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		response := scanner.Text()
		fmt.Println("Server response:", response)

		// Check if the response is a Merkle Root
		if len(response) == 64 {
			fmt.Println("Received Merkle Root:", response)
		} else if strings.HasPrefix(response, "Transaction verification result:") {
			parts := strings.Split(response, ":")
			result := strings.TrimSpace(parts[1])

			if result == "true" {
				fmt.Println("Transaction exists in the Merkle tree.")
			} else {
				fmt.Println("Transaction does not exist in the Merkle tree.")
			}
		}
	}
}

func buildMerkleTree(conn net.Conn) {
	// Send the BUILD_MERKLE_TREE command
	sendCommand(conn, "BUILD_MERKLE_TREE")
}

// ...


func sendCommand(conn net.Conn, command string) {
	_, err := conn.Write([]byte(command + "\n"))
	if err != nil {
		fmt.Println("Error sending command:", err)
	}
}

func createBlock(conn net.Conn, transactions []string) {
	// Send the ADD_BLOCK command with a list of transactions
	command := fmt.Sprintf("ADD_BLOCK|%s", strings.Join(transactions, ","))
	sendCommand(conn, command)
}

func queryTransaction(conn net.Conn, transactionData string, merkleRoot string) {
	// Send the QUERY_TRANSACTION command with Merkle Root and Transaction Data
	command := fmt.Sprintf("QUERY_TRANSACTION|%s|%s", merkleRoot, transactionData)
	sendCommand(conn, command)
}


func main() {
	// Use flag to get the server port from the command line
	port := flag.Int("port", 8080, "Server port")
	flag.Parse()

	serverAddress := fmt.Sprintf("localhost:%d", *port)

	// Connect to the server
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
		if !scanner.Scan() {
			fmt.Println("Error reading from standard input:", scanner.Err())
			break
		}
		command := scanner.Text()

		if strings.HasPrefix(command, "CREATE_BLOCK") {
			// Example: CREATE_BLOCK|transaction1,transaction2,transaction3
			parts := strings.Split(command, "|")
			transactions := strings.Split(parts[1], ",")
			createBlock(conn, transactions)
		} else if strings.HasPrefix(command, "QUERY_TRANSACTION") {
			// Example: QUERY_TRANSACTION|merkleRoot|transaction_data
			parts := strings.Split(command, "|")
			merkleRoot := parts[1]
			transactionData := parts[2]
			queryTransaction(conn, transactionData, merkleRoot)
		} else if command == "BUILD_MERKLE_TREE" {
			// Request to build the Merkle Tree
			buildMerkleTree(conn)
		} else {
			sendCommand(conn, command)
		}
	}

	fmt.Println("Connection closed.")
}
