// Package state deals with game states and transitions.
package state

import (
	"crypto/rand"
	"fmt"
	"log"
	"sort"

	"golang.org/x/crypto/blake2b"
)

// MinPlayers is the minimum number of players in a game.
const MinPlayers = 2

// A Message carries the input data that causes state transitions.
type Message struct {
	From int
	Data []byte
}

// A Fn represents a state. Calling Fn transitions into the next state.
type Fn func(m *Machine) Fn

// Machine represents the game state machine.
type Machine struct {
	NPlayers int
	WhoAmI   int
	In       <-chan Message
	Out      chan<- Message

	nextPlayer int // players are numbered 1..N
	err        error

	ss    [][8]byte
	hs    [][32]byte
	order []int

	debug bool
}

// Run runs the state machine until it terminates, returning a non-nil error if
// execution failed.
func (m *Machine) Run() error {
	if m.NPlayers < MinPlayers {
		return fmt.Errorf("too few players: got %d, want %d or more", m.NPlayers, MinPlayers)
	}
	if m.WhoAmI < 1 || m.WhoAmI > m.NPlayers {
		return fmt.Errorf("invalid player identification: %d not in range [1,%d]", m.WhoAmI, m.NPlayers)
	}
	if m.In == nil {
		return fmt.Errorf("m.In is nil")
	}
	if m.Out == nil {
		return fmt.Errorf("m.Out is nil")
	}
	for state := stateGameplayOrder1PublishH; state != nil; {
		state = state(m)
	}
	return m.Err()
}

// Err returns the error associated with this machine. A non-nil error means the
// machine terminated in a failure state.
func (m *Machine) Err() error {
	return m.err
}

// Fail returns the nil Fn, denoting termination, and sets the error associated
// with this machine.
func (m *Machine) Fail(err error) Fn {
	m.err = err
	return nil
}

func (m *Machine) logf(format string, args ...interface{}) {
	if m.debug {
		log.Printf("Player #%d: %s", m.WhoAmI, fmt.Sprintf(format, args...))
	}
}

func stateGameplayOrder1PublishH(m *Machine) Fn {
	m.nextPlayer = m.nextPlayer%m.NPlayers + 1

	if m.WhoAmI == m.nextPlayer {
		var s [8]byte
		if _, err := rand.Read(s[:]); err != nil {
			return m.Fail(err)
		}
		h := blake2b.Sum256(s[:])
		m.ss = append(m.ss, s)
		m.hs = append(m.hs, h)
		m.logf("Out <- %x", h)
		m.Out <- Message{
			From: m.WhoAmI,
			Data: h[:],
		}
	}

	msg, ok := <-m.In
	if !ok {
		return m.Fail(fmt.Errorf("expected more messages"))
	}
	if msg.From != m.nextPlayer {
		return m.Fail(fmt.Errorf("message from unexpected player: got %v, want %v", msg.From, m.nextPlayer))
	}
	m.logf("In -> %x", msg.Data)
	var got [blake2b.Size256]byte
	copy(got[:], msg.Data)

	if msg.From == m.WhoAmI {
		want := m.hs[m.WhoAmI-1]
		if want != got {
			return m.Fail(fmt.Errorf("corrupted message: want %x, got %x", want, got))
		}
	} else {
		m.hs = append(m.hs, got)
	}

	if len(m.hs) == m.NPlayers {
		return stateGameplayOrder2PublishS
	}

	return stateGameplayOrder1PublishH
}

func stateGameplayOrder2PublishS(m *Machine) Fn {
	m.nextPlayer = m.nextPlayer%m.NPlayers + 1

	if m.WhoAmI == m.nextPlayer {
		s := m.ss[m.WhoAmI-1]
		m.logf("Out <- %x", s)
		m.Out <- Message{
			From: m.WhoAmI,
			Data: s[:],
		}
	}

	msg, ok := <-m.In
	if !ok {
		return m.Fail(fmt.Errorf("expected more messages"))
	}
	if msg.From != m.nextPlayer {
		return m.Fail(fmt.Errorf("message from unexpected player: got %v, want %v", msg.From, m.nextPlayer))
	}
	m.logf("In -> %x", msg.Data)
	var got [8]byte
	copy(got[:], msg.Data)

	if msg.From == m.WhoAmI {
		want := m.ss[m.WhoAmI-1]
		if want != got {
			return m.Fail(fmt.Errorf("corrupted message: want %x, got %x", want, got))
		}
	} else {
		if got, want := blake2b.Sum256(msg.Data), m.hs[msg.From-1]; got != want {
			return m.Fail(fmt.Errorf("hash of %x does not match: got %x, want %x", msg.Data, got, want))
		}
		m.ss = append(m.ss, got)
	}

	if len(m.ss) == m.NPlayers {
		return stateGameplayOrder3Compute
	}

	return stateGameplayOrder2PublishS
}

func stateGameplayOrder3Compute(m *Machine) Fn {
	t := blake2b.Sum256(xor(m.ss...))
	m.order = m.order[:0]
	for i := 1; i <= m.NPlayers; i++ {
		m.order = append(m.order, i)
	}
	sort.SliceStable(m.order, func(i int, j int) bool {
		return string(t[i*8:i*8+8]) < string(t[j*8:j*8+8])
	})
	m.logf("gameplay order: %v", m.order)
	return stateShuffleTiles
}

func xor(ss ...[8]byte) []byte {
	if len(ss) == 0 {
		return nil
	}
	out := ss[0][:]
	for i, s := range ss[1:] {
		out[i] ^= s[i]
	}
	return out
}

func stateShuffleTiles(m *Machine) Fn {
	return m.Fail(fmt.Errorf("state not implemented yet: waiting for Player #%d to shuffle tiles", m.order[0]))
}
