package model

type LastHop struct {
	ID       string
	IP       string
	Port     string
	NodeType string
}

func NewLastHop(id string, ip string, port string, nodetype string) *LastHop {
	return &LastHop{
		ID:       id,
		IP:       ip,
		Port:     port,
		NodeType: nodetype,
	}
}
