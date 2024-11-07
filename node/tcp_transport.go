package node

import (
	"errors"
	"fmt"
	"net"
)

// check interface implementation
var _ Transport = &TCPTransport{}

type TCPTransportOptions struct {
	ListenAddress string
	Decoder       Decoder
	HandshakeFunc HandshakeFunc
	OnNode        func(n *TCPNode) error
}

type TCPTransport struct {
	TCPTransportOptions
	listener net.Listener
	msg      chan MSG
}

func NewTCPTransport(options TCPTransportOptions) *TCPTransport {
	return &TCPTransport{
		TCPTransportOptions: options,
		msg:                 make(chan MSG),
	}
}

func (t *TCPTransport) Addr() string {
	return t.ListenAddress
}
func (t *TCPTransport) Close() error {
	return t.Close()
}
func (t *TCPTransport) Consume() <-chan MSG {
	return t.msg
}

func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	go t.handleConn(conn, true)

	return nil
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddress)
	if err != nil {
		fmt.Println("TCP listen error: ", err)
		return err
	}

	fmt.Println("TCP listen: ", t.listener.Addr())

	go t.acceptLoop()

	return nil
}

func (t *TCPTransport) acceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			fmt.Println("TCP accept error: ", err)
		}

		go t.handleConn(conn, false)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn, isOut bool) {
	var err error

	node := NewTCPNode(conn, isOut)

	if err = t.HandshakeFunc(node); err != nil {
		fmt.Println("TCP handshake error:", err)
		return
	}

	if t.OnNode != nil {
		if err = t.OnNode(node); err != nil {
			fmt.Println("TCP OnNode error:", err)
			return
		}
	}

	//////////////////////////
	for {
		return
	}
}
