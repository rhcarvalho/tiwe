package main

import (
	"context"
	"fmt"
	"log"

	"github.com/rhcarvalho/tiwe/crypto/rand"
	"github.com/rhcarvalho/tiwe/router"
)

const maxMessages = 6

type Player struct {
	Name string
	router.Router

	players []string

	cancel context.CancelFunc
	done   chan struct{}

	recv func() router.Message
	send func(v interface{}) error

	msgs []*router.Message
}

func (p *Player) Start(ctx context.Context, players []string) error {
	p.players = players

	c, err := p.Register(ctx)
	if err != nil {
		return err
	}
	recv := func(v interface{}) error {
		msg, err := c.Recv(ctx)
		v = msg
		p.logf("<-- %v", v)
		return err
	}
	p.send = func(v interface{}) (err error) {
		msg := v.(router.Message)
		msg.From = p.Name
		p.logf("--> %v", msg)
		err = c.Send(ctx, msg)
		return err
	}

	// ch := make(chan *Game)

	ctx, p.cancel = context.WithCancel(ctx)
	p.done = make(chan struct{})
	go func() {
		defer close(p.done)
		_, _ = recv, p.send
		p.loop(ctx, nil)
	}()
	return nil
}

func (p *Player) Init() {
	n := len(p.players)
	data := make([]byte, n)
	for i, v := range rand.Perm(n) {
		data[i] = byte(v)
	}
	p.send(router.Message{ID: 1, Data: data})
}

func (p *Player) Stop(ctx context.Context) error {
	p.cancel()
	select {
	case <-p.Done():
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *Player) Done() <-chan struct{} {
	return p.done
}

func (p *Player) loop(ctx context.Context, update <-chan *Game) {
	startState := waitFor(p.players[0])
	for state := startState; state != nil; {
		select {
		case <-ctx.Done():
			return
		case game := <-update:
			state = state(game)
		}
	}
	// for {
	// 	var msg Message
	// 	err := recv(&msg)
	// 	if err != nil {
	// 		p.logf("receive: Decode: %v", err)
	// 		return
	// 	}
	// 	if msg.ID == maxMessages {
	// 		return
	// 	}
	// 	p.msgs = append(p.msgs, &msg)
	// 	if p.isMyTurn() {
	// 		// play()
	// 		msg.ID++
	// 		rand.Shuffle(len(msg.Data), func(i int, j int) {
	// 			msg.Data[i], msg.Data[j] = msg.Data[j], msg.Data[i]
	// 		})
	// 		func() {
	// 			err := send(msg)
	// 			if err != nil {
	// 				p.logf("receive: Broadcast: %v", err)
	// 				return
	// 			}
	// 		}()
	// 	}
	// }
}

type Game struct {
	player       string
	initialOrder []string

	Messages []*router.Message
	Err      error

	order []string
}

func (g *Game) NextPlayer() string {
	order := g.initialOrder
	// if g.order != nil {
	// 	order = g.order
	// }
	if len(g.Messages) == 0 {
		return order[0]
	}
	last := g.Messages[len(g.Messages)-1].From
	lastIndex := SliceIndex(len(order), func(i int) bool { return order[i] == last })
	next := order[(lastIndex+1)%len(order)]
	return next
}

func (g *Game) myTurn() bool {
	order := g.initialOrder
	// if g.order != nil {
	// 	order = g.order
	// }
	if len(g.Messages) == 0 {
		return len(order) > 0 && g.player == order[0]
	}
	last := g.Messages[len(g.Messages)-1].From
	lastIndex := SliceIndex(len(order), func(i int) bool { return order[i] == last })
	next := order[(lastIndex+1)%len(order)]
	return g.player == next
}

func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

type stateFn func(g *Game) stateFn

func waitFor(player string) stateFn {
	return func(g *Game) stateFn {
		if len(g.Messages) == 0 {
			g.Err = fmt.Errorf("no messages")
			return nil
		}
		last := g.Messages[len(g.Messages)-1].From
		if last == player {
			if g.NextPlayer() == g.player {
				//g.Play()
			}
			return waitFor(g.NextPlayer())
		}
		g.Err = fmt.Errorf("waiting for %q, got message from %q", player, last)
		return nil
	}
}

// func (p *Player) isMyTurn() bool {
// 	if len(p.msgs) == 0 {
// 		return false
// 	}
// 	last := p.msgs[len(p.msgs)-1]
// 	return last.From == p.After
// }

func (p *Player) logf(format string, args ...interface{}) {
	log.Printf("%-25s %v",
		fmt.Sprintf("Player{Name: %q}:", p.Name),
		fmt.Sprintf(format, args...))
}
