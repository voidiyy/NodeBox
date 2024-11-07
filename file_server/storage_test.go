package fileSrv

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestCASPath(t *testing.T) {
	key := "lalala"

	expectedFullPath := "df2ef/a060e/335f9/7628c/a39c9/fef54/69ab3/cb837/df2efa060e335f97628ca39c9fef5469ab3cb837"
	expectedFirstBlock := "df2ef"
	path := CASPath(key)

	if path.FullPath() != expectedFullPath {
		t.Errorf("full path mismatch need: %s got: %s", expectedFullPath, path.FullPath())
	}

	if path.FirstPathBlock() != expectedFirstBlock {
		t.Errorf("block path mismatch need: %s got: %s", expectedFirstBlock, path.FirstPathBlock())
	}
}

func TestCreateStorage(t *testing.T) {
	s := newStorage()
	id := "randID"

	key := fmt.Sprintf("storage")
	data := []byte("some data lalalalalalaal")

	n, err := s.Write(key, id, bytes.NewReader(data))
	if err != nil {
		t.Error(err)
	}

	if n != int64(len(data)) {
		t.Errorf("error writeng excected data len: %d -- writed data: %d", len(data), n)
	}

}

func TestFullStorage(t *testing.T) {
	s := newStorage()
	id := "randID"

	defer testClear(t, s)

	for i := 0; i < 50; i++ {
		key := fmt.Sprintf("storage__%d", i)
		data := []byte("some data lalalalalalaal")

		if _, err := s.Write(key, id, bytes.NewReader(data)); err != nil {
			t.Error(err)
		}

		if ok := s.Has(key, id); !ok {
			t.Errorf("Has func not get a file")
		}

		n, r, err := s.Read(key, id)
		if err != nil {
			t.Error(err)
		}

		fmt.Printf("s.Read: read %d bytes\n", n)

		b, _ := io.ReadAll(r)
		if string(b) != string(data) {
			t.Errorf("expected data: %d  readed data: %d", data, b)
		}

		if err := s.DeleteFile(key, id); err != nil {
			t.Error(err)
		}

		if ok := s.Has(key, id); ok == true {
			t.Errorf("s.Has: file %s must be deleted", key)
		}

	}

}

func newStorage() *Storage {
	opts := StorageOpts{
		PathFunc: CASPath,
	}

	return NewStorage(opts)
}

func testClear(t *testing.T, s *Storage) {
	err := s.ClearFull()
	if err != nil {
		t.Error(err)
	}
}
