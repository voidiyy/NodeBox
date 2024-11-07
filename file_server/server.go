package fileSrv

import (
	"NodeBox/node"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"sync"
)

const (
	IncomingMessage = 0x1
	IncomingStream  = 0x2
)

type FileServerOpts struct {
	ID                string
	EncKey            []byte
	StorageRoot       string
	PathTransformFunc PathFunc
	Transport         *node.TCPTransport
	BootstrapNodes    []string
}

type FileServer struct {
	FileServerOpts

	nodeLock sync.Mutex
	nodes    map[string]*node.TCPNode

	Store    *Storage
	quitChan chan struct{}
}

func NewFileServer(opts FileServerOpts) (*FileServer, error) {
	s := StorageOpts{
		Root:     opts.StorageRoot,
		PathFunc: opts.PathTransformFunc,
	}

	fs := &FileServer{
		FileServerOpts: opts,
		nodeLock:       sync.Mutex{},
		nodes:          make(map[string]*node.TCPNode),
		Store:          NewStorage(s),
		quitChan:       make(chan struct{}),
	}

	return fs, nil
}

func (fs *FileServer) Start() {
	fmt.Printf("node [%s] starting . . .\n", fs.Transport.Addr())

	if err := fs.Transport.ListenAndAccept(); err != nil {
		fmt.Println(err)
	}

	for {
		fs.bootstrapNetwork()

		fs.loop()
	}
}

func (fs *FileServer) broadcast(msg *Message) error {
	buff := new(bytes.Buffer)

	if err := gob.NewEncoder(buff).Encode(msg); err != nil {
		return err
	}

	for _, n := range fs.nodes {
		if err := n.Send([]byte{IncomingMessage}); err != nil {
			return err
		}

		if err := n.Send(buff.Bytes()); err != nil {
			return err
		}
	}

	return nil
}

func (fs *FileServer) Stop() {
	close(fs.quitChan)
}

func (fs *FileServer) OnNode(n *node.TCPNode) error {
	fs.nodeLock.Lock()
	defer fs.nodeLock.Unlock()

	fs.nodes[n.RemoteAddr().String()] = n

	return nil
}

func (fs *FileServer) loop() {
	defer func() {
		fmt.Printf("file server stoped via error or user exit")
		fs.Transport.Close()
	}()

	for {
		select {
		case m := <-fs.Transport.Consume():
			var msg Message

			if err := gob.NewDecoder(bytes.NewReader(m.Payload)).Decode(&msg); err != nil {
				log.Println("[loop] decoding message error: ", err)
			}

			if err := fs.handleMessage(m.Sender, &msg); err != nil {
				log.Println("[loop] handling message error: ", err)
			}
		case <-fs.quitChan:
			log.Println("quit channel")
			return
		}
	}
}

func (fs *FileServer) bootstrapNetwork() {
	for _, addr := range fs.BootstrapNodes {
		if len(addr) == 0 {
			continue
		}

		go func(addr string) {
			err := fs.Transport.Dial(addr)
			if err != nil {
				log.Println("dial error with remote node ", addr)
			}

			fmt.Printf("node [%s] dial connection with remote node %s\n", fs.Transport.Addr(), addr)
		}(addr)
	}
}

func init() {
	gob.Register(MessageStoreFile{})
	gob.Register(MessageGetFile{})
}
