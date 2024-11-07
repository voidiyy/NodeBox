package node

import "net"

// Node implementation file - tcp_node.go
type Node interface {
	net.Conn
	Send([]byte) error
	CloseStream()
}

// Transport implementation file - tcp_transport.go
type Transport interface {
	Addr() string
	Dial(string) error
	ListenAndAccept() error
	Close() error
	Consume() <-chan MSG
}
