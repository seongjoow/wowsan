// network.go
package network

// import (
// 	"fmt"
// )

// // Network 구조체는 브로커 노드들의 네트워크를 나타냅니다.
// type Network struct {
// 	Brokers map[string]*Broker // 브로커 노드만을 관리하는 맵
// }

// // Broker 구조체는 브로커의 상태와 속성을 나타냅니다.
// type Broker struct {
// 	ID          string
// 	Publishers  map[string]*Publisher
// 	Subscribers map[string]*Subscriber
// }

// // NewNetwork 함수는 새로운 네트워크를 생성합니다.
// func NewNetwork() *Network {
// 	return &Network{
// 		Brokers: make(map[string]*Broker),
// 	}
// }

// // NewBroker 함수는 새로운 브로커를 생성하고 네트워크에 추가합니다.
// func (n *Network) NewBroker(id string) *Broker {
// 	broker := &Broker{
// 		ID:          id,
// 		Publishers:  make(map[string]*Publisher),
// 		Subscribers: make(map[string]*Subscriber),
// 	}
// 	n.Brokers[id] = broker
// 	fmt.Printf("Broker %s created and added to network\n", id)
// 	return broker
// }

// // 여기에 다른 네트워크 관련 함수들을 추가...
