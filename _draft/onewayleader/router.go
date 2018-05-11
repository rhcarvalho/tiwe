package main

import (
	"context"
	"net"
	"sync"
)

type Router interface {
	Register(ctx context.Context) (net.Conn, error)
	Broadcast(ctx context.Context, data []byte) error
}

type testRouter struct {
	mu    sync.RWMutex
	conns []net.Conn
}

func (r *testRouter) Register(ctx context.Context) (net.Conn, error) {
	c1, c2 := net.Pipe()
	r.mu.Lock()
	r.conns = append(r.conns, c1)
	r.mu.Unlock()
	return c2, nil
}

func (r *testRouter) Broadcast(ctx context.Context, data []byte) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, conn := range r.conns {
		select {
		case <-ctx.Done():
		default:
			if t, ok := ctx.Deadline(); ok {
				err := conn.SetWriteDeadline(t)
				if err != nil {
					return err
				}
			}
			_, err := conn.Write(data)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *testRouter) len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.conns)
}
