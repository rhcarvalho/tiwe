package main

import (
	"encoding/gob"
	"net"
	"sync"
	"testing"
	"time"
)

func Test(t *testing.T) {
	in, out := net.Pipe()

	var wg sync.WaitGroup
	wg.Add(1)
	enc := gob.NewEncoder(out)
	go func() {
		start := time.Now()
		defer wg.Done()
		dec := gob.NewDecoder(in)
		for {
			if time.Since(start) > 8*time.Second {
				break
			}
			var msg Message
			err := dec.Decode(&msg)
			if err != nil {
				t.Fatal(err)
			}
			t.Log(msg)
			msg.ID++
			go func() {
				err := enc.Encode(msg)
				if err != nil {
					t.Fatal(err)
				}
			}()
		}
	}()

	err := enc.Encode(Message{From: "Alice", ID: 1, Data: []byte("abc")})
	if err != nil {
		t.Fatal(err)
	}

	wg.Wait()
}
