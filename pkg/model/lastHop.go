package model

type LastHop struct {
	Id       string
	Ip       string
	Port     string
	NodeType string
}

func NewLastHop(id string, ip string, port string, nodetype string) *LastHop {
	return &LastHop{
		Id:       id,
		Ip:       ip,
		Port:     port,
		NodeType: nodetype,
	}
}
