package game

import (
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"
)

type Message struct {
	From string
}

func Play(name string, ch chan Message) {
	select {
	case ch <- Message{From: name}:
	case <-time.After(time.Duration(rand.Intn(100)) * time.Millisecond):
	}
	select {
	case m := <-ch:
		log.Printf("From: %s, To: %s", m.From, name)
	case <-time.After(time.Duration(rand.Intn(100)) * time.Millisecond):
	}
}

func TestGame(t *testing.T) {
	ch := make(chan Message)
	var wg sync.WaitGroup
	names := []string{"Alice", "Bob", "Carol"}
	wg.Add(len(names))
	for _, name := range names {
		name := name
		go func() {
			defer wg.Done()
			Play(name, ch)
		}()
	}
	wg.Wait()
}
