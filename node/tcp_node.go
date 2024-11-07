package node

import (
	"net"
	"sync"
)

var _ Node = &TCPNode{}

type TCPNode struct {
	net.Conn

	// if connection outbound => true
	// if connection inbound => false
	isOut bool
	wg    *sync.WaitGroup
}

func NewTCPNode(conn net.Conn, isOut bool) *TCPNode {
	return &TCPNode{
		Conn:  conn,
		isOut: isOut,
		wg:    &sync.WaitGroup{},
	}
}

func (n *TCPNode) Send(b []byte) error {
	_, err := n.Conn.Write(b)
	return err
}

func (n *TCPNode) CloseStream() {
	n.wg.Done()
}
