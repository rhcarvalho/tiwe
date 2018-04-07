// Package commutative implements a commutative encryption scheme such that a
// message encrypted multiple times with different keys can be decrypted with
// the same keys but in any order.
package commutative

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"golang.org/x/crypto/blake2b"
)

const (
	keySize   = 16 // AES-128 key size.
	NonceSize = keySize
	Size256   = 32 // 256-bit cryptographic hash size in bytes.
)

// A Message holds bytes to be encrypted and/or decrypted.
type Message struct {
	Bytes  []byte
	Nonces map[[Size256]byte][NonceSize]byte
}

// NewMessage returns a new Message with the given bytes. The new Message takes
// ownership of b, and the caller should not use b after this call.
func NewMessage(b []byte) *Message {
	return &Message{
		Bytes:  b,
		Nonces: make(map[[Size256]byte][NonceSize]byte),
	}
}

// A Key is used to encrypt and decrypt messages.
type Key struct {
	sum   [Size256]byte
	block cipher.Block
}

// GenerateKey returns a new random key.
func GenerateKey() *Key {
	key := make([]byte, keySize)
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	return &Key{
		sum:   blake2b.Sum256(key),
		block: block,
	}
}

// Encrypt uses k to encrypt the message in place. It returns the message to
// make it convenient to chain calls.
func (k *Key) Encrypt(m *Message) *Message {
	if _, ok := m.Nonces[k.sum]; ok {
		panic("attempt to encrypt twice with the same key")
	}
	var iv [NonceSize]byte
	if _, err := rand.Read(iv[:]); err != nil {
		panic(err)
	}
	m.Nonces[k.sum] = iv
	stream := cipher.NewCTR(k.block, iv[:])
	stream.XORKeyStream(m.Bytes, m.Bytes)
	return m
}

// Decrypt uses k to decrypt the message in place. It returns the message to
// make it convenient to chain calls.
func (k *Key) Decrypt(m *Message) *Message {
	iv, ok := m.Nonces[k.sum]
	if !ok {
		panic("attempt to decrypt before encrypt")
	}
	delete(m.Nonces, k.sum)
	stream := cipher.NewCTR(k.block, iv[:])
	stream.XORKeyStream(m.Bytes, m.Bytes)
	return m
}

// CanDecrypt reports whether key can decrypt message.
func (k *Key) CanDecrypt(m *Message) bool {
	_, ok := m.Nonces[k.sum]
	return ok
}
