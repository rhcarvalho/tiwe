// Package state deals with game states and transitions.
package state

import (
	"fmt"
	"log"

	"github.com/rhcarvalho/tiwe/crypto/commutative"
	"github.com/rhcarvalho/tiwe/crypto/rand"
)

// MinPlayers is the minimum number of players in a game.
const MinPlayers = 2

// A Message carries the input data that causes state transitions.
type Message struct {
	From          int
	EncryptedData *commutative.Message
}

// A Fn represents a state. Calling Fn with a message transitions into
// a new state.
type Fn func(m *Machine, msg Message) Fn

// Machine represents the game state machine.
type Machine struct {
	NPlayers int
	WhoAmI   int
	In       <-chan Message
	Out      chan<- Message

	nextPlayer int // players are numbered 1..N
	err        error

	key *commutative.Key

	debug bool
}

// Run runs the state machine until it terminates, returning a non-nil error if
// execution failed.
func (m *Machine) Run() error {
	for state := m.Start(); state != nil; {
		msg, ok := <-m.In
		if !ok {
			state = m.Fail(fmt.Errorf("expected more messages"))
			continue
		}
		state = state(m, msg)
	}
	return m.Err()
}

// Start returns the initial state function of the machine.
func (m *Machine) Start() Fn {
	if m.NPlayers < MinPlayers {
		return m.Fail(fmt.Errorf("too few players: got %d, want %d or more", m.NPlayers, MinPlayers))
	}
	m.nextPlayer = 1

	// In a game with N players, the first player (in some implicit order)
	// creates a slice of bytes with incremental values from 0 to N-1. Next,
	// it shuffles and encrypts the bytes with a random key using
	// commutative encryption.
	if m.WhoAmI == m.nextPlayer {
		p := make([]byte, m.NPlayers)
		for i := range p {
			p[i] = byte(i)
		}
		e := commutative.NewMessage(p)
		m.key = shuffleEncrypt(e)
		m.logf("Out <- % x", e.Bytes)
		m.Out <- Message{
			From:          m.WhoAmI,
			EncryptedData: e,
		}
	}
	return stateDecideGameplayOrder
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

func stateDecideGameplayOrder(m *Machine, msg Message) Fn {
	if msg.From != m.nextPlayer {
		return m.Fail(fmt.Errorf("message from unexpected player: got %v, want %v", msg.From, m.nextPlayer))
	}
	if len(msg.EncryptedData.Bytes) != m.NPlayers {
		return m.Fail(fmt.Errorf("encrypted data field: got %d bytes, want %d", len(msg.EncryptedData.Bytes), m.NPlayers))
	}
	m.logf("In -> % x", msg.EncryptedData.Bytes)
	if m.key != nil && !m.key.CanDecrypt(msg.EncryptedData) {
		// This error implies that the m.nextPlayer player somehow
		// removed the encrytion nonce used with m.key, indicating a
		// misbehavior (e.g. trying to manipulate the game state).
		return m.Fail(fmt.Errorf("bad encrypted data: message cannot be decrypted"))
	}

	m.nextPlayer = m.nextPlayer%m.NPlayers + 1
	if m.nextPlayer == 1 {
		return stateDecideGameplayOrderRevealSecret
	}

	if m.WhoAmI == m.nextPlayer {
		m.key = shuffleEncrypt(msg.EncryptedData)
		m.logf("Out <- % x", msg.EncryptedData.Bytes)
		m.Out <- Message{
			From:          m.WhoAmI,
			EncryptedData: msg.EncryptedData,
		}
	}

	return stateDecideGameplayOrder
}

func stateDecideGameplayOrderRevealSecret(m *Machine, msg Message) Fn {
	return m.Fail(fmt.Errorf("state not implemented yet"))
}

// shuffleEncrypt shuffles the bytes in m and encrypts them with a random key.
// It mutates m and returns the key.
func shuffleEncrypt(m *commutative.Message) *commutative.Key {
	rand.Shuffle(len(m.Bytes), func(i int, j int) {
		m.Bytes[i], m.Bytes[j] = m.Bytes[j], m.Bytes[i]
	})
	key := commutative.GenerateKey()
	key.Encrypt(m)
	return key
}
