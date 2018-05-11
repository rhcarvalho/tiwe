package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"net"

	"github.com/rhcarvalho/tiwe/crypto/rand"
)

type Player struct {
	Name string
	Router
}

func (p *Player) Start(ctx context.Context) { // TODO: g Game
	conn, err := p.Register(ctx)
	if err != nil {
		p.logf("Start: %v", err)
		panic(err)
	}
	p.receive(ctx, conn)
}

func (p *Player) Broadcast(ctx context.Context, msg Message) error {
	msg.From = p.Name
	data := mustMarshal(msg)
	p.logf("--> %v <%v>", msg, data)
	return p.Router.Broadcast(ctx, data)
}

func (p *Player) receive(ctx context.Context, conn net.Conn) {
	if t, ok := ctx.Deadline(); ok {
		conn.SetReadDeadline(t)
	}
	// NOTE: doesn't work, why? Because once the decoder learns about a new
	// type, namely Message, it doesn't expect to see the type defined
	// again. But we create a new encoder for every message, and thus the
	// encoder resends the definition of Message over and over, and a single
	// receive decoder does not handle that behavior.
	// dec := gob.NewDecoder(&ObservableConn{conn})
	for {
		dec := gob.NewDecoder(conn) // NOTE: works, why?! See above.
		var msg Message
		err := dec.Decode(&msg)
		if err != nil {
			p.logf("receive: Decode: %v", err)
			return
		}
		p.logf("<-- %v", msg)
		if true { // TODO: isMyTurn() {
			// play()
			msg.ID++
			rand.Shuffle(len(msg.Data), func(i int, j int) {
				msg.Data[i], msg.Data[j] = msg.Data[j], msg.Data[i]
			})
			go func() {
				err := p.Broadcast(ctx, msg)
				if err != nil {
					p.logf("receive: Broadcast: %v", err)
					return
				}
			}()
		}
	}
}

func (p *Player) logf(format string, args ...interface{}) {
	log.Printf("%-25s %v",
		fmt.Sprintf("Player{Name: %q}:", p.Name),
		fmt.Sprintf(format, args...))
}

func marshal(v interface{}) ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func mustMarshal(v interface{}) []byte {
	b, err := marshal(v)
	if err != nil {
		panic("mustMarshal: " + err.Error())
	}
	return b
}
