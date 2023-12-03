package network

import (
	"fmt"
	"net"
)

func handleConnection(conn net.Conn, nodeID int) {
	defer conn.Close()
	fmt.Printf("Node %d connected: %s\n", nodeID, conn.RemoteAddr().String())

	// Xử lý kết nối ở đây

}

func StartNode(nodeID, port int) {
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
		go handleConnection(conn, nodeID)
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
	