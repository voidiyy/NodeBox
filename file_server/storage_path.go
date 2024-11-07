package fileSrv

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
)

// Path type represent path to file, and folder tree
type Path struct {
	Path string
	File string
}

// CASPath generate dir tree with depth 5, and hash values for file ands dirs
func CASPath(key string) Path {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 10
	pathLen := 1

	paths := make([]string, pathLen)

	for i := 0; i < pathLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashStr[from:to]
	}

	return Path{
		Path: strings.Join(paths, "/"),
		File: hashStr,
	}
}

// DefaultPathFunc generate simple root and file path
var DefaultPathFunc = func(key string) Path {
	return Path{
		Path: key,
		File: key,
	}
}

// FirstPathBlock return first dir from dir tree
func (p *Path) FirstPathBlock() string {
	paths := strings.Split(p.Path, "/")

	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

// FullPath return full file path
func (p *Path) FullPath() string {
	return fmt.Sprintf("%s/%s", p.Path, p.File)
}
