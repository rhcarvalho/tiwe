package router

import (
	"context"
	"encoding/gob"
)

// A Router broadcasts messages to all registered parties.
type Router interface {
	Register(ctx context.Context) (Conn, error)
}

// A Conn allows sending messages to and receiving messages from a Router.
type Conn interface {
	Send(ctx context.Context, m Message) error
	Recv(ctx context.Context) (Message, error)
}

// A Message is a container to transfer data among parties.
type Message struct {
	From string
	ID   int
	Data []byte
}

// gobConn implements Conn by marshaling messages using the encoding/gob format.
type gobConn struct {
	enc *gob.Encoder
	dec *gob.Decoder
}

// Send implements Conn.
func (c *gobConn) Send(ctx context.Context, m Message) error {
	return c.enc.Encode(m)
}

// Recv implements Conn.
func (c *gobConn) Recv(ctx context.Context) (Message, error) {
	var m Message
	return m, c.dec.Decode(&m)
}
