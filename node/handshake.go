package node

type HandshakeFunc func(node *TCPNode) error

func NILHandshake(node *TCPNode) error {
	return nil
}
