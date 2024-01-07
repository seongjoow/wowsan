package model

type LastHop struct {
	ID   string
	IP   string
	Port string
}

func NewLastHop(id string, ip string, port string) *LastHop {
	return &LastHop{
		ID:   id,
		IP:   ip,
		Port: port,
	}
}
