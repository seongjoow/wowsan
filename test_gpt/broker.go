package main

// func (s *broker) SendMessage(ctx context.Context, in *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {

// 	return &pb.SendMessageResponse{
// 		Message: in.Message,
// 		}, nil
// }

// // interface
// type Broker interface {
// 	RegisterPublisher(publisher *model.Publisher)
// 	RegisterSubscriber(subscriber *model.Subscriber)
// 	ProcessMessages()
// }

// // ProcessMessages 함수는 발행자로부터 메시지를 수신하고 구독자에게 전달하는 로직을 포함합니다.
// func (b *broker) ProcessMessages() {

// 	for {
// 		msg, err := b.Node.ReceiveMessage()
// 		if err != nil {
// 			fmt.Println("Error receiving message:", err)
// 			continue
// 		}

// 		for _, sub := range b.Subscribers {
// 			err := sub.Node.SendMessage(msg)
// 			if err != nil {
// 				fmt.Printf("Error sending message to subscriber %s: %v\n", sub.ID, err)
// 			}
// 		}
// 	}
// }

// // 여기에 추가적인 브로커 관련 로직을 구현...
