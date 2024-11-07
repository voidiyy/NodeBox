package fileSrv

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

// CopyDecrypt accept cipher key and two data sources
// create init vector and read to it from source, where it decrypted
// crete CTR and pass all into copyStream func()
func CopyDecrypt(key []byte, src io.Reader, dst io.Writer) (int, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, err
	}

	iv := make([]byte, block.BlockSize())

	if _, err = src.Read(iv); err != nil {
		return 0, err
	}

	stream := cipher.NewCTR(block, iv)
	return copyStream(stream, block.BlockSize(), src, dst)
}

// CopyEncrypt accept cipher key and two data sources
// create init vector and read and encrypt to it from source
// then write encrypted data into dst writer
// crete CTR and pass all into copyStream func()
func CopyEncrypt(key []byte, src io.Reader, dst io.Writer) (int, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, err
	}

	iv := make([]byte, block.BlockSize())

	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return 0, err
	}

	if _, err = dst.Write(iv); err != nil {
		return 0, err
	}

	stream := cipher.NewCTR(block, iv)
	return copyStream(stream, block.BlockSize(), src, dst)
}

// copyStream func accept encrypted/decrypted data
// read from src and write to dst until EOF
// return wrote amount of bytes
func copyStream(stream cipher.Stream, blockSize int, src io.Reader, dst io.Writer) (int, error) {
	var (
		buff = make([]byte, 32*1024)
		bs   = blockSize
	)

	for {
		n, err := src.Read(buff)
		if n > 0 {
			stream.XORKeyStream(buff, buff[:n])
			nn, err := dst.Write(buff[:n])
			if err != nil {
				return 0, err
			}
			bs += nn
		}

		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return 0, err
		}
	}

	return bs, nil
}

func generateID() string {
	buff := make([]byte, 12)
	_, _ = io.ReadFull(rand.Reader, buff)
	return hex.EncodeToString(buff)
}

func hashKey(key string) string {
	hash := md5.Sum([]byte(key))
	return hex.EncodeToString(hash[:])
}

func NewEncKey() []byte {
	buff := make([]byte, 32)
	_, _ = io.ReadFull(rand.Reader, buff)
	return buff
}
