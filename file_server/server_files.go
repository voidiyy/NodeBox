package fileSrv

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"time"
)

func (fs *FileServer) GetFile(key string) (io.Reader, error) {
	if fs.Store.Has(key, fs.ID) {
		fmt.Printf("node [%s], get file %s\n", fs.Transport.Addr(), key)
		n, r, err := fs.Store.Read(key, fs.ID)
		if r == nil || err != nil {
			return nil, fmt.Errorf("local read error or file not found")
		}
		fmt.Printf("node [%s] file successfully read with %d bytes\n", fs.Transport.Addr(), n)
		return r, err
	}

	fmt.Printf("node [%s] file not exist locally, trying to fetch from network\n", fs.Transport.Addr())

	msg := &Message{
		Payload: &MessageGetFile{
			Key: hashKey(key),
			ID:  fs.ID,
		},
	}

	if err := fs.broadcast(msg); err != nil {
		return nil, err
	}

	time.Sleep(1 * time.Second)

	if len(fs.nodes) == 0 {
		return nil, fmt.Errorf("no nodes available for file fetch")
	}

	for _, n := range fs.nodes {

		if n.Conn == nil {
			fmt.Printf("node [%s] remote node is nil %s\n", fs.Transport.Addr(), n.RemoteAddr())
			continue
		}

		var fileSize int64
		if err := binary.Read(n, binary.LittleEndian, &fileSize); err != nil {
			log.Println("error reading file size:", err)
			continue
		}

		nn, err := fs.Store.WriteDecrypt(fs.EncKey, key, fs.ID, io.LimitReader(n, fileSize))
		if err != nil {
			return nil, err
		}

		fmt.Printf("node [%s] received %d bytes from remote node [%s]\n", fs.Transport.Addr(), nn, n.RemoteAddr())
		n.CloseStream()
	}

	_, rr, err := fs.Store.Read(key, fs.ID)
	if rr == nil || err != nil {
		return nil, fmt.Errorf("file could not be read after fetch")
	}
	fmt.Println("before return")
	return rr, err
}

func (fs *FileServer) StoreFile(key string, r io.Reader) error {
	var (
		fileBuff = new(bytes.Buffer)
		tee      = io.TeeReader(r, fileBuff)
	)

	size, err := fs.Store.Write(key, fs.ID, tee)
	if err != nil {
		return err
	}

	msg := &Message{
		Payload: &MessageStoreFile{
			Size: size + 16,
			Key:  hashKey(key),
			ID:   fs.ID,
		},
	}

	if err = fs.broadcast(msg); err != nil {
		return err
	}

	time.Sleep(1 * time.Second)

	nodes := []io.Writer{}

	for _, n := range fs.nodes {
		nodes = append(nodes, n)
	}

	mw := io.MultiWriter(nodes...)
	mw.Write([]byte{IncomingStream})

	n, err := CopyEncrypt(fs.EncKey, fileBuff, mw)
	if err != nil {
		return err
	}

	fmt.Printf("node [%s] received and writed %d bytes to local storage\n", fs.Transport.Addr(), n)
	return nil
}
