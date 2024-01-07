package model

type Node struct {
	ID   string
	IP   string
	Port int
}

func NewNode(id string, ip string, port int) *Node {
	return &Node{
		ID:   id,
		IP:   ip,
		Port: port,
	}
}
