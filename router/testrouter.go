package router

import (
	"context"
	"encoding/gob"
	"io"
	"sync"
	"time"

	"github.com/rhcarvalho/tiwe/crypto/rand"
)

// A TestRouter implements an in-memory Router, independent of a network.
type TestRouter struct {
	// Latency defines the parameters to simulate network latency. The
	// observed latency follows a normal distribution with the given mean
	// and standard deviation.
	Latency struct {
		Mean   time.Duration
		StdDev time.Duration
	}

	mu   sync.RWMutex
	encs []*gob.Encoder
}

// NewTestRouter return a new TestRouter with the given network latency
// parameters.
func NewTestRouter(mean, sd time.Duration) *TestRouter {
	r := TestRouter{}
	r.Latency.Mean = mean
	r.Latency.StdDev = sd
	return &r
}

// Register registers a new party. It returns a Conn that can be used to
// broadcast and receive messages.
func (r *TestRouter) Register(ctx context.Context) (Conn, error) {
	in1, out1 := io.Pipe()
	in2, out2 := io.Pipe()
	enc1 := gob.NewEncoder(out1)
	dec1 := gob.NewDecoder(in1)
	enc2 := gob.NewEncoder(out2)
	dec2 := gob.NewDecoder(in2)
	r.mu.Lock()
	r.encs = append(r.encs, enc2)
	r.mu.Unlock()
	go func() {
		for {
			var msg Message
			err := dec1.Decode(&msg)
			if err != nil {
				panic(err)
			}
			r.mu.RLock()
			for _, enc := range r.encs {
				r.simulateNetworkLatency()
				err := enc.Encode(msg)
				if err != nil {
					panic(err)
				}
			}
			r.mu.RUnlock()
		}
	}()
	return &gobConn{
		enc: enc1,
		dec: dec2,
	}, nil
}

// simulateNetworkLatency simulates network latency by sleeping for a random
// duration with normal distribution.
func (r *TestRouter) simulateNetworkLatency() {
	d := time.Duration(rand.NormFloat64()*float64(r.Latency.StdDev) + float64(r.Latency.Mean))
	time.Sleep(d)
}
