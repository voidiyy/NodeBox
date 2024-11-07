package fileSrv

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Message struct {
	Payload any
}

type MessageStoreFile struct {
	Size int64
	Key  string
	ID   string
}

type MessageGetFile struct {
	Key string
	ID  string
}

func (fs *FileServer) handleMessage(from string, msg *Message) error {
	switch v := msg.Payload.(type) {
	case MessageStoreFile:
		return fs.handleMessageStore(from, v)
	case MessageGetFile:
		return fs.handleMessageGet(from, v)
	}
	return nil
}

func (fs *FileServer) handleMessageGet(from string, msg MessageGetFile) error {
	if !fs.Store.Has(msg.Key, msg.ID) {
		return fmt.Errorf("node [%s] not found %s file locally...\n", fs.Transport.Addr(), msg.Key)
	}

	fmt.Printf("node [%s] serving file %s over the network...\n", fs.Transport.Addr(), msg.Key)

	fileSize, fileReader, err := fs.Store.Read(msg.Key, msg.ID)
	if err != nil {
		return err
	}

	if rc, ok := fileReader.(io.ReadCloser); ok {
		fmt.Printf("node [%s] closing readeCloser\n")
		defer rc.Close()
	}

	noda, ok := fs.nodes[from]
	if !ok {
		return fmt.Errorf("node [%s] not found remote node %s in nodeList\n", fs.Transport.Addr(), from)
	}

	err = noda.Send([]byte{IncomingStream})
	if err != nil {
		return fmt.Errorf("node [%s] error send first byte to remote node %s\n", fs.Transport.Addr(), noda.RemoteAddr())
	}

	err = binary.Write(noda, binary.LittleEndian, fileSize)
	if err != nil {
		return fmt.Errorf("node [%s] error send file size to remote node %s\n", fs.Transport.Addr(), noda.RemoteAddr())
	}

	n, er := io.Copy(noda, fileReader)
	if er != nil {
		return fmt.Errorf("node [%s] error copy data to remote node %s\n", fs.Transport.Addr(), noda.RemoteAddr())
	}

	fmt.Printf("node [%s] successfully send %d bytes to remote node %s\n", fs.Transport.Addr(), n, noda.RemoteAddr())

	return nil
}

func (fs *FileServer) handleMessageStore(from string, msg MessageStoreFile) error {
	noda, ok := fs.nodes[from]
	if !ok {
		return fmt.Errorf("node [%s] not found remote node %s in nodeList\n", fs.Transport.Addr(), from)
	}

	n, err := fs.Store.Write(msg.Key, msg.ID, io.LimitReader(noda, msg.Size))
	if err != nil {
		return fmt.Errorf("node [%s] error write data to remote node %s", fs.Transport.Addr(), noda.RemoteAddr())
	}

	fmt.Printf("node [%s] successfully stored %d bytes from remote node %s\n", fs.Transport.Addr(), n, noda.RemoteAddr())

	noda.CloseStream()

	return nil
}
