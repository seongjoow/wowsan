package transport_pool

import (
	"fmt"
	"sync"

	"google.golang.org/grpc"
)

type ConnectionPool struct {
	pool map[string][]*grpc.ClientConn // Key 為 broker 地址
	mu   sync.Mutex
}

func NewConnectionPool() *ConnectionPool {
	return &ConnectionPool{
		pool: make(map[string][]*grpc.ClientConn),
	}
}

// Add connection to pool
func (cp *ConnectionPool) AddConnection(broker string) error {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	conn, err := grpc.Dial(broker, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to connect to broker %s: %v", broker, err)
	}
	cp.pool[broker] = append(cp.pool[broker], conn)
	return nil
}

func (cp *ConnectionPool) GetConnection(broker string) (*grpc.ClientConn, error) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	if conns, exists := cp.pool[broker]; exists && len(conns) > 0 {
		conn := conns[0]
		cp.pool[broker] = conns[1:] // pop the first connection
		return conn, nil
	}
	return nil, fmt.Errorf("no available connection for broker: %s", broker)
}

// After using the connection, return it to the pool
func (cp *ConnectionPool) ReturnConnection(broker string, conn *grpc.ClientConn) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	cp.pool[broker] = append(cp.pool[broker], conn)
}

// Close all connections
func (cp *ConnectionPool) CloseAll() {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	for broker, conns := range cp.pool {
		for _, conn := range conns {
			conn.Close()
		}
		delete(cp.pool, broker)
	}
}
