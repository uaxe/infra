package zhash

import (
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"hash"
	"sync"
)

// BcryptHash
func BcryptHash(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

// BcryptCheck
func BcryptCheck(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type Hasher struct {
	pool sync.Pool
}

func NewHasher(h hash.Hash) *Hasher {
	hasher := &Hasher{pool: sync.Pool{New: func() any { return h }}}
	return hasher
}

func (h *Hasher) Get(key string) string {
	hasher := h.pool.Get().(hash.Hash)
	defer h.pool.Put(hasher)
	hasher.Write([]byte(key))
	hashstr := hex.EncodeToString(hasher.Sum(nil))
	hasher.Reset()
	return hashstr
}
