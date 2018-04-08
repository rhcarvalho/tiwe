package state

import (
	"testing"

	"golang.org/x/crypto/blake2b"
)

func TestStateMachine(t *testing.T) {
	in := make(chan Message, 1)
	out := make(chan Message, 1)
	m := Machine{
		NPlayers: 3,
		WhoAmI:   1,
		In:       in,
		Out:      out,

		debug: true,
	}
	go func() {
		msg1 := <-out
		in <- msg1

		s2 := []byte("foo")
		h2 := blake2b.Sum256(s2)
		in <- Message{
			From: 2,
			Data: h2[:],
		}

		s3 := []byte("bar")
		h3 := blake2b.Sum256(s3)
		in <- Message{
			From: 3,
			Data: h3[:],
		}

		in <- <-out
		in <- Message{2, s2}
		in <- Message{3, s3}

		in <- Message{}
	}()
	err := m.Run()
	if err != nil {
		t.Fatal(err)
	}
}
