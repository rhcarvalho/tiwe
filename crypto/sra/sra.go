// Package sra implements a commutative encryption scheme, as described in
// Shamir, Rivest and Adleman's (SRA) Mental Poker.
//
// The encryption is such that a plain text message may be encrypted by one or
// more keys yielding a cipher text that can be decrypted by the same keys in
// any order.
//
// Citation
//
// Shamir, A., Rivest, R. L., & Adleman, L. M. (1981). Mental poker. In The
// mathematical gardner (pp. 37-43). Springer, Boston, MA.
//
//  http://people.csail.mit.edu/rivest/ShamirRivestAdleman-MentalPoker.pdf
package sra

import (
	"crypto/rand"
	"io"
	"math/big"
)

// SRA encryption requires a large prime number. This package uses a fixed
// documented prime instead of generating one.
// This is the 1024-bit MODP Group with 160-bit Prime Order Subgroup from
// RFC 5114, section 2.1.
//
//  https://tools.ietf.org/html/rfc5114#section-2.1
const (
	primeHex = "B10B8F96A080E01DDE92DE5EAE5D54EC52C99FBCFB06A3C69A6A9DCA52D23B616073E28675A23D189838EF1E2EE652C013ECB4AEA906112324975C3CD49B83BFACCBDD7D90C4BD7098488E9C219A73724EFFD6FAE5644738FAA31A4FF55BCCC0A151AF5F0DC8B4BD45BF37DF365C1A65E68CFDA76D4DA708DF1FB2BC2E4A4371"
	// generatorHex = "A4D1CBD5C3FD34126765A442EFB99905F8104DD258AC507FD6406CFF14266D31266FEA1E5C41564B777E690F5504F213160217B4B01B886A5E91547F9E2749F4D7FBD7D3B9A92EE1909D0D2263F80A76A6A24C087A091F531DBF0A0169B6A28AD662A4D18E73AFA32D779D5918D08BC8858F4DCEF97C2A24855E6EEB22B3B2E5"
	// qHex = "F518AA8781A8DF278ABA4E7D64B7CB9D49462353"
)

// minBitLen is the minimum number of bits required for generated encryption and
// decryption exponents.
const minBitLen = 160

// Reusable big.Int values, allocated only once globally.
var (
	bigOne = big.NewInt(1)
	// defaultN is a large prime used in the encryption scheme.
	defaultN = fromHex(primeHex)
	// defaultTotient is the value of Euler's totient function applied to
	// defaultN. Since defaultN is prime, the totient is trivially N-1.
	defaultTotient = new(big.Int).Sub(defaultN, bigOne)
	// defaultMaxK is the upper limit for the encryption exponent K.
	// Limiting K makes encryption faster, and, consequently, reduces the
	// total time for encryption+decryption too.
	defaultMaxK = new(big.Int).Lsh(bigOne, minBitLen)
)

func fromHex(hex string) *big.Int {
	n, ok := new(big.Int).SetString(hex, 16)
	if !ok {
		panic("bad hex number: " + hex)
	}
	return n
}

// A Key represents an SRA key.
type Key struct {
	// N is a large prime.
	N *big.Int
	// K is the encryption exponent.
	K *big.Int
	// L is the decryption exponent.
	L *big.Int
}

// GenerateKey generates a Key using the given random source (e.g.,
// crypto/rand.Reader).
func GenerateKey(random io.Reader) *Key {
	key := &Key{N: defaultN}
	var err error
	g := new(big.Int)
start:
	key.K, err = rand.Int(random, defaultMaxK)
	if err != nil {
		panic("sra: cannot generate random number: " + err.Error())
	}
	// Set K's highest bit to ensure that it has minBitLen
	// significant bits.
	key.K.SetBit(key.K, minBitLen-1, 1)
	// Set K's lowest bit to ensure it is odd. In order to have
	// GCD(K, Phi(N)) == 1, K must be odd, because Phi(N) is even
	// (it is a prime minus one).
	key.K.SetBit(key.K, 0, 1)
	g.GCD(nil, nil, key.K, defaultTotient)
	if g.Cmp(bigOne) != 0 {
		goto start
	}
	// Repurpose g to avoid an allocation.
	key.L = g.ModInverse(key.K, defaultTotient)
	if key.L.BitLen() < minBitLen {
		goto start
	}
	return key
}

// Encrypt encrypts plaintext.
func (k *Key) Encrypt(plaintext []byte) []byte {
	m := new(big.Int).SetBytes(plaintext)
	if m.Cmp(k.N) >= 0 {
		panic("plaintext is too long")
	}
	c := new(big.Int).Exp(m, k.K, k.N)
	return c.Bytes()
}

// Decrypt decrypts ciphertext.
func (k *Key) Decrypt(ciphertext []byte) []byte {
	c := new(big.Int).SetBytes(ciphertext)
	if c.Cmp(k.N) >= 0 {
		panic("ciphertext is invalid")
	}
	m := new(big.Int).Exp(c, k.L, k.N)
	return m.Bytes()
}
