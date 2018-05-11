package main

import (
	"context"
	"log"
	"time"
)

type Player struct {
	Name string
	Router
	pendingGames []Game
	currentGames []Game
}

func (p *Player) Start() {
	parent := context.Background()
	for {
		log.Printf("%s: looping!", p.Name)
		ctx, cancel := context.WithTimeout(parent, 10*time.Second)
		msg, err := p.Receive(ctx)
		if err != nil {
			log.Printf("%s: Receive: %v", p.Name, err)
			continue
		}
		p.HandleMessage(msg)
		cancel()
	}
}

func (p *Player) HandleMessage(m *Message) {
	log.Printf("%s: HandleMessage: %v", p.Name, m)
	switch m.Type {
	case MessageTypeRequestToPlay:
		p.Broadcast(context.TODO(), &Message{
			Type: MessageTypeAcceptToPlay,
		})
	case MessageTypeAcceptToPlay:
		log.Printf("%s: HandleMessage: let's play!", p.Name)
	default:
		log.Printf("%s: HandleMessage: unknown type %v", p.Name, m.Type)
	}
}

func (p *Player) NewGame(who ...interface{}) {
	p.Broadcast(context.TODO(), &Message{
		Type: MessageTypeRequestToPlay,
	})
}

type Game struct {
	Router
}

type MessageType uint8

const (
	MessageTypeUndefined MessageType = iota
	MessageTypeRequestToPlay
	MessageTypeAcceptToPlay
)

type Message struct {
	Type MessageType
	// From string
	// To   []string
}

type Router interface {
	Broadcast(context.Context, *Message) error
	Receive(context.Context) (*Message, error)
}

type testRouter struct {
	out map[string]chan<- *Message
	in  <-chan *Message
}

func newTestRouter(key string, m map[string]chan *Message) Router {
	router := &testRouter{
		out: make(map[string]chan<- *Message),
	}
	for k, v := range m {
		if k == key {
			router.in = v
			continue
		}
		router.out[k] = v
	}
	if router.in == nil {
		panic("newTestRouter: missing channel for " + key)
	}
	return router
}

func (r *testRouter) Broadcast(ctx context.Context, msg *Message) error {
	var sg SyncGroup
	for _, ch := range r.out {
		ch := ch
		sg.Go(func() {
			select {
			case ch <- msg:
			case <-ctx.Done():
			}
		})
	}
	sg.Wait()
	return ctx.Err()
}

func (r *testRouter) Receive(ctx context.Context) (*Message, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case msg := <-r.in:
		return msg, nil
	}
}

func main() {
	cs := map[string]chan *Message{
		"alice": make(chan *Message),
		"bob":   make(chan *Message),
		"carol": make(chan *Message),
	}
	alice := &Player{
		Name:   "alice",
		Router: newTestRouter("alice", cs),
	}
	bob := &Player{
		Name:   "bob",
		Router: newTestRouter("bob", cs),
	}
	carol := &Player{
		Name:   "carol",
		Router: newTestRouter("carol", cs),
	}
	var sg SyncGroup
	sg.Go(alice.Start)
	sg.Go(bob.Start)
	sg.Go(carol.Start)
	alice.NewGame("bob", "carol")
	sg.Wait()
}
