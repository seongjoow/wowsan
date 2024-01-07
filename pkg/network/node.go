// node.go
package network

// import (
// 	"fmt"
// 	"net"
// 	"time"
// )

// type Node struct {
// 	ID   string
// 	Addr string
// 	Port int
// 	Conn net.Conn
// }

// func NewNode(id, addr string, port int) *Node {
// 	return &Node{
// 		ID:   id,
// 		Addr: addr,
// 		Port: port,
// 	}
// }

// // Connect 함수는 다른 노드에 연결을 시도합니다.
// func (n *Node) Connect() error {
// 	var err error
// 	n.Conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", n.Addr, n.Port))
// 	if err != nil {
// 		return fmt.Errorf("failed to connect to node %s: %v", n.ID, err)
// 	}
// 	fmt.Printf("Connected to node %s\n", n.ID)
// 	return nil
// }

// // Disconnect 함수는 다른 노드와의 연결을 해제합니다.
// func (n *Node) Disconnect() {
// 	if n.Conn != nil {
// 		n.Conn.Close()
// 		fmt.Printf("Disconnected from node %s\n", n.ID)
// 	}
// }

// // SendMessage 함수는 연결된 노드에 메시지를 전송합니다.
// func (n *Node) SendMessage(msg string) error {
// 	_, err := n.Conn.Write([]byte(msg))
// 	if err != nil {
// 		return fmt.Errorf("failed to send message to node %s: %v", n.ID, err)
// 	}
// 	return nil
// }

// // ReceiveMessage 함수는 노드로부터 메시지를 수신합니다.
// func (n *Node) ReceiveMessage() (string, error) {
// 	buffer := make([]byte, 1024)
// 	n.Conn.SetReadDeadline(time.Now().Add(5 * time.Second))
// 	readLen, err := n.Conn.Read(buffer)
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(buffer[:readLen]), nil
// }
