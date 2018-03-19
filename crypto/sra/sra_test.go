package sra

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"
)

func TestGenerateKey(t *testing.T) {
	start := time.Now()
	const (
		minKeys     = 100
		maxKeys     = 1000
		maxDuration = 500 * time.Millisecond
	)
	seen := make(map[string]struct{}, maxKeys)
	for i := 0; i < maxKeys; i++ {
		// Limit test time to maxDuration after seeing minKeys.
		if i > minKeys && time.Since(start) > maxDuration {
			break
		}
		key := GenerateKey(rand.Reader)
		kstr := key.K.String()
		if _, ok := seen[kstr]; ok {
			t.Fatalf("K already seen: %s (after %d iterations)", kstr, i)
		}
		seen[kstr] = struct{}{}
		if key.N != defaultN {
			t.Errorf("N = %v, want %v", key.N, defaultN)
		}
		if key.K.BitLen() < minBitLen {
			t.Errorf("K[=%v].BitLen() = %v, want >= %v", key.K, key.K.BitLen(), minBitLen)
		}
		if key.L.BitLen() < minBitLen {
			t.Errorf("L[=%v].BitLen() = %v, want >= %v", key.L, key.L.BitLen(), minBitLen)
		}
	}
}

func TestEncrytDecrypt(t *testing.T) {
	plaintext := []byte("secret message")
	key := GenerateKey(rand.Reader)
	ciphertext := key.Encrypt(plaintext)
	if string(ciphertext) == string(plaintext) {
		t.Errorf("key.Encrypt: got plaintext")
	}
	recovered := key.Decrypt(ciphertext)
	if got, want := string(recovered), string(plaintext); got != want {
		t.Errorf("key.Decrypt: got %q, want %q", got, want)
	}
}

func TestCommutative(t *testing.T) {
	plaintext := []byte("secret message")
	key1 := GenerateKey(rand.Reader)
	key2 := GenerateKey(rand.Reader)
	key3 := GenerateKey(rand.Reader)
	perms := [][]*Key{
		{key1, key2, key3},
		{key1, key3, key2},
		{key2, key1, key3},
		{key2, key3, key1},
		{key3, key1, key2},
		{key3, key2, key1},
	}
	for _, permEncryption := range perms {
		for _, permDecryption := range perms {
			buf := dup(plaintext)
			for _, key := range permEncryption {
				buf = key.Encrypt(buf)
			}
			for _, key := range permDecryption {
				buf = key.Decrypt(buf)
			}
			if !eq(buf, plaintext) {
				t.Fatalf("not commutative: got %q, want %q", buf, plaintext)
			}
		}
	}
}

func dup(src []byte) (dst []byte) {
	dst = make([]byte, len(src))
	copy(dst, src)
	return dst
}

func eq(a, b []byte) bool {
	return string(a) == string(b)
}

// BenchmarkExponentBitLen shows how smaller exponents lead to faster modular
// exponentiation. The purpose is to justify limiting the encryption exponent K
// in generated keys.
func BenchmarkExponentBitLen(b *testing.B) {
	max := new(big.Int)
	z := new(big.Int)
	for _, bitLen := range []uint{64, 128, 160, 1024, 2048} {
		max.Lsh(bigOne, bitLen)
		e, err := rand.Int(rand.Reader, max)
		if err != nil {
			b.Fatal(err)
		}
		// Set highest bit to ensure bitLen significant bits.
		e.SetBit(e, int(bitLen)-1, 1)
		// Set lowest bit to ensure it is odd.
		e.SetBit(e, 0, 1)
		b.Run(strconv.FormatUint(uint64(bitLen), 10), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				z.Exp(big.NewInt(int64(i)), e, defaultN)
			}
		})
	}
}

// TODO: check quadratic residue

func TestQuadraticResidues(t *testing.T) {
	var r rune
	for i, found := 0, 0; found < 106; i++ {
		res := big.Jacobi(big.NewInt(int64(i)), defaultN)
		r = ' '
		if res == 1 {
			r = '*'
			found++
		}
		fmt.Printf("%d\t%s\t%d\n", i, string(r), res)
	}
}
