package pkg

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"

	"golang.org/x/crypto/sha3"
)

type HashFun struct {
	hash.Hash
}

func GetHashFun(hashName string) HashFun {
	switch hashName {
	case "SHA1":
		return HashFun{sha1.New()}
	case "MD5":
		return HashFun{md5.New()}
	case "SHA256":
		return HashFun{sha256.New()}
	case "SHA3_256":
		return HashFun{sha3.New256()}
	case "SHA512_256":
		return HashFun{sha512.New512_256()}
	case "SHA512":
		return HashFun{sha512.New()}
	default:
		return HashFun{sha1.New()}
	}
}

func (h *HashFun) GetHash(data []byte) string {
	h.Reset()
	h.Write(data)
	return string(h.Sum(nil))
}
