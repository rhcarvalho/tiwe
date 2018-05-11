package main

import (
	"crypto/hmac"
	"fmt"
	"hash"

	"golang.org/x/crypto/blake2b"
)

func main() {
	key := []byte("secret")
	msg := []byte("Hello, playground")

	hm := hmac.New(func() hash.Hash {
		h, _ := blake2b.New256(nil)
		return h
	}, key)
	hm.Write(msg)
	fmt.Println(hm.Sum(nil))

	h, _ := blake2b.New256(key)
	h.Write(msg)
	fmt.Println(h.Sum(nil))
}
