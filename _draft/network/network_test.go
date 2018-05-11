package network

import (
	"context"
	"net"
	"sync"
	"testing"
)

type Game struct{}

func (g *Game) Over() bool {
	return false
}

func (g *Game) WaitForTurn(ctx context.Context) {

}

func (g *Game) Play() {

}

type Player struct {
	name   string
	dialer Dialer

	mu       sync.RWMutex
	listener net.Listener
	ready    <-chan struct{}
}

func NewPlayer(name string) *Player {
	return &Player{
		name: name,
	}
}

func (p *Player) Accept(l net.Listener) {
	p.mu.Lock()
	p.listener = l
	ch := make(chan struct{})
	p.ready = ch
	p.mu.Unlock()
	close(ch)
	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		_ = conn
	}
}

// Returns a channel that is closed when the player is ready to play.
func (p *Player) Ready() <-chan struct{} {
	return p.ready
}

func (p *Player) Name() string {
	return p.name
}

func (p *Player) Addr() net.Addr {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.listener != nil {
		return p.listener.Addr()
	}
	return nil
}

func (p *Player) Start(args ...interface{}) (*Game, error) {
	return nil, nil
}

func (p *Player) NewRequestToPlay(addrs ...net.Addr) *RequestToPlay {
	return &RequestToPlay{
		From: p,
		To:   addrs,
	}
}

type RequestToPlay struct {
	From *Player
	To   []net.Addr
}

func (rtp *RequestToPlay) Send(ctx context.Context) ([]net.Addr, error) {
	for _, addr := range rtp.To {
		conn, err := rtp.From.dialer.DialContext(ctx, addr.Network(), addr.String())
		if err != nil {
			continue
		}
		_ = conn
	}
	return nil, nil
}

func TestDecideGameplayOrder(t *testing.T) {
	alice := NewPlayer("Alice")
	bob := NewPlayer("Bob")
	carol := NewPlayer("Carol")

	dialer := &testNetwork{}
	alice.dialer = dialer
	bob.dialer = dialer
	carol.dialer = dialer

	var wg sync.WaitGroup
	listen := func(p *Player) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p.Accept(dialer.Listen("test", p.Name()))
		}()
	}
	listen(alice)
	listen(bob)
	listen(carol)

	<-alice.Ready()
	<-bob.Ready()
	<-carol.Ready()

	rtp := alice.NewRequestToPlay(bob.Addr(), carol.Addr())
	addrs, err := rtp.Send(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	game, err := alice.Start(context.TODO(), addrs)
	if err != nil {
		t.Fatal(err)
	}
	for !game.Over() {
		game.WaitForTurn(context.TODO())
		game.Play()
	}
}

type Dialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

// testNetwork implements fake network connections.
type testNetwork struct {
	sync.Once

	sync.Mutex
	conns map[addr]<-chan net.Conn
}

type addr struct {
	network, address string
}

func (tn *testNetwork) init() {
	tn.conns = make(map[addr]<-chan net.Conn)
}

func (tn *testNetwork) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	tn.Do(tn.init)
	tn.Lock()
	ch := tn.conns[addr{network, address}]
	tn.Unlock()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case conn := <-ch:
		return conn, nil
	}
}

func (tn *testNetwork) Listen(network, address string) net.Listener {
	tn.Do(tn.init)
	ch := make(chan net.Conn)
	tn.Lock()
	tn.conns[addr{network, address}] = ch
	tn.Unlock()
	return &testListener{
		ch: ch,
	}
}

type testListener struct {
	ch chan<- net.Conn
}

// Accept waits for and returns the next connection to the listener.
func (tl *testListener) Accept() (net.Conn, error) {
	c1, c2 := net.Pipe()
	tl.ch <- c1
	return c2, nil
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (tl *testListener) Close() error {
	return nil
}

// Addr returns the listener's network address.
func (tl *testListener) Addr() net.Addr {
	return nil
}
