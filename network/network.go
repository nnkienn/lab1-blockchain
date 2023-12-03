package network

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"github.com/nnkienn/lab1-blockchain/server/blockchain"


)

var nodes = []string{"127.0.0.1:3001", "127.0.0.1:3002", "127.0.0.1:3003"}
var chainMutex sync.Mutex

func handleConnection(conn net.Conn, nodeID int, chain *blockchain.Blockchain) {
	defer conn.Close()
	fmt.Printf("Node %d connected: %s\n", nodeID, conn.RemoteAddr().String())

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		receivedMsg := scanner.Text()
		fmt.Printf("Received message from client: %s\n", receivedMsg)

		switch receivedMsg {
		case "printchain":
			chainMutex.Lock()
			for _, block := range chain.GetBlocks() {
				blockInfo := fmt.Sprintf("Timestamp: %d\nPrev. hash: %x\n", block.Timestamp, block.PrevBlockHash)
				for _, transaction := range block.Transactions {
					blockInfo += fmt.Sprintf("Transaction: %s\n", string(transaction.Data))
				}
				blockInfo += fmt.Sprintf("Hash: %x\n", block.Hash)

				_, err := fmt.Fprintf(conn, "%s\n", blockInfo)
				if err != nil {
					fmt.Println(err)
					chainMutex.Unlock()
					return
				}
			}
			chainMutex.Unlock()

		case "hello":
			responseMsg := "Hello from node!"
			_, err := fmt.Fprintf(conn, "%s\n", responseMsg)
			if err != nil {
				fmt.Println(err)
				return
			}

		default:
			fmt.Fprintln(conn, "Invalid command")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func StartNode(nodeID, port int, chain *blockchain.Blockchain) {
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		fmt.Printf("Node %d error listening: %s\n", nodeID, err.Error())
		return
	}
	defer listener.Close()
	fmt.Printf("Node %d started. Listening on 127.0.0.1:%d\n", nodeID, port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Node %d error accepting connection: %s\n", nodeID, err.Error())
			return
		}
		go handleConnection(conn, nodeID, chain)
	}
}

func ConnectNodes(sourceNode int, host string, port int) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		fmt.Printf("Node %d error connecting to Node at %s:%d: %s\n", sourceNode, host, port, err.Error())
		return
	}
	defer conn.Close()
	fmt.Printf("Node %d connected to Node at %s:%d\n", sourceNode, host, port)
}
