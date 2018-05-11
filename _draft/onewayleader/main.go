package main

import (
	"context"
	"log"
	"time"
)

func main() {
	var router testRouter

	alice := &Player{
		Name:   "alice",
		Router: &router,
	}
	bob := &Player{
		Name:   "bob",
		Router: &router,
	}
	carol := &Player{
		Name:   "carol",
		Router: &router,
	}

	ctx := context.Background()
	go alice.Start(ctx)
	go bob.Start(ctx)
	go carol.Start(ctx)

	// Wait until all players are ready.
	for i := uint(0); router.len() < 3; i++ {
		d := time.Duration(1<<(3*i)) * time.Millisecond
		if d > 10*time.Second {
			log.Fatal("players not ready")
		}
		time.Sleep(d)
	}

	alice.Broadcast(ctx, Message{ID: 1, Data: []byte{1, 2, 3}})

	time.Sleep(10 * time.Second)
}

type Message struct {
	From string
	ID   int
	Data []byte
}
