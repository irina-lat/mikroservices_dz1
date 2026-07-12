package grpc

import (
	"google.golang.org/grpc"
)

// Client представляет gRPC клиент
type Client struct {
	conn *grpc.ClientConn
}

// NewClient создаёт новый gRPC клиент
func NewClient(conn *grpc.ClientConn) *Client {
	return &Client{conn: conn}
}

// Close закрывает gRPC соединение
func (c *Client) Close() error {
	return c.conn.Close()
}