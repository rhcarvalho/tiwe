package rand

import (
	"testing"
	"time"
)

func TestSourceInt63(t *testing.T) {
	const (
		maxDuration   = 100 * time.Millisecond
		maxIterations = 100000
	)
	var last int64
	start := time.Now()
	for i := 1; time.Since(start) < maxDuration && i <= maxIterations; i++ {
		v := source.Int63()
		if v < 0 {
			t.Fatalf("v = %d, after %d iterations", v, i)
		}
		if v == last {
			t.Fatalf("v = %d == last seen, after %d iterations", v, i)
		}
		last = v
	}
}

func TestSourceUint64(t *testing.T) {
	const (
		maxDuration   = 100 * time.Millisecond
		maxIterations = 100000
	)
	var last uint64
	var bits uint64
	start := time.Now()
	for i := 1; time.Since(start) < maxDuration && i <= maxIterations; i++ {
		v := source.Uint64()
		if v == last {
			t.Fatalf("v = %d == last seen, after %d iterations", v, i)
		}
		last = v
		bits |= v
	}
	if bits != 1<<64-1 {
		t.Errorf("bits = %064b, want all bits set", bits)
	}
}
