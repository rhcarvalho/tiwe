package rand

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"math/rand"
)

var source rand.Source64 = cryptoRandSource{}

// cryptoRandSource implements math/rand.Source64 backed by crypto/rand.
type cryptoRandSource struct{}

func (cryptoRandSource) Seed(int64) {}

func (r cryptoRandSource) Int63() int64 {
	// &^ (1 << 63) clears the sign bit
	return int64(r.Uint64() &^ (1 << 63))
}

func (cryptoRandSource) Uint64() uint64 {
	var b [8]byte
	_, err := cryptorand.Read(b[:])
	if err != nil {
		panic(err)
	}
	return binary.BigEndian.Uint64(b[:])
}
