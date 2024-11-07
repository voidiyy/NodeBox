package fileSrv

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

// default value for root folder
const defaultRoot = "nodeStorage"

type PathFunc func(string) Path

// StorageOpts represent name of root folder and func for creating file tree
type StorageOpts struct {
	Root     string
	PathFunc PathFunc
}

type Storage struct {
	StorageOpts
}

// NewStorage create a storage entity with specified or default values
func NewStorage(opts StorageOpts) *Storage {
	if opts.PathFunc == nil {
		opts.PathFunc = DefaultPathFunc
	}
	if opts.Root == "" {
		opts.Root = defaultRoot
	}

	return &Storage{
		opts,
	}
}

// Has func check if file with provide key & id if exist
func (s *Storage) Has(key, id string) bool {
	path := s.PathFunc(key)
	fullPath := fmt.Sprintf("%s/%s/%s", s.Root, id, path.Path)

	_, err := os.Stat(fullPath)
	return !errors.Is(err, os.ErrNotExist)
}

// ClearFull func remove root dir
func (s *Storage) ClearFull() error {
	return os.RemoveAll(s.Root)
}

// DeleteFile func delete file from storage by key & id
func (s *Storage) DeleteFile(key, id string) error {
	path := s.PathFunc(key)

	defer func() {
		log.Printf("file [%s] deleted from disk", path.File)
	}()

	fullPath := fmt.Sprintf("%s/%s/%s", s.Root, id, path.FirstPathBlock())

	return os.RemoveAll(fullPath)
}

// Read func start s.readStream
func (s *Storage) Read(key, id string) (int64, io.ReadCloser, error) {
	return s.readStream(key, id)
}

// Write func start functions chain
func (s *Storage) Write(key, id string, r io.Reader) (int64, error) {
	return s.writeStream(key, id, r)
}

// WriteDecrypt func creating a file and
// write to it decrypted data via func in crypto.go file
func (s *Storage) WriteDecrypt(enc []byte, key, id string, r io.Reader) (int64, error) {
	f, err := s.openFileWrite(key, id)
	if err != nil {
		return 0, err
	}
	n, er := CopyDecrypt(enc, r, f)
	if er != nil {
		return 0, er
	}

	return int64(n), nil
}

// openFileWrite func crete dirs and file by key & id
func (s *Storage) openFileWrite(key, id string) (*os.File, error) {
	path := s.PathFunc(key)

	fullPathDir := fmt.Sprintf("%s/%s/%s", s.Root, id, path.Path)

	if err := os.MkdirAll(fullPathDir, os.ModePerm); err != nil {
		return nil, err
	}

	fullPathFile := fmt.Sprintf("%s/%s/%s", s.Root, id, path.FullPath())

	return os.Create(fullPathFile)
}

// writeStream func copy from reader to newly created file
func (s *Storage) writeStream(key, id string, r io.Reader) (int64, error) {
	f, err := s.openFileWrite(key, id)
	if err != nil {
		return 0, err
	}

	return io.Copy(f, r)
}

// readStream func open file via key & id and return file size and read closer
func (s *Storage) readStream(key, id string) (int64, io.ReadCloser, error) {
	path := s.PathFunc(key)
	fullPath := fmt.Sprintf("%s/%s/%s", s.Root, id, path.FullPath())

	f, err := os.Open(fullPath)
	if err != nil {
		return 0, nil, err
	}

	stat, err := f.Stat()
	if err != nil {
		return 0, nil, err
	}

	return stat.Size(), f, nil
}
