package state

import (
	"testing"

	"github.com/rhcarvalho/tiwe/crypto/commutative"
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

		msg2 := deepCopy(msg1)
		msg2.From = 2
		shuffleEncrypt(msg2.EncryptedData)
		in <- msg2

		msg3 := deepCopy(msg2)
		msg3.From = 3
		shuffleEncrypt(msg3.EncryptedData)
		in <- msg3

		in <- Message{}
	}()
	err := m.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func deepCopy(msg Message) Message {
	b := make([]byte, len(msg.EncryptedData.Bytes))
	copy(b, msg.EncryptedData.Bytes)
	new := Message{
		From:          msg.From,
		EncryptedData: commutative.NewMessage(b),
	}
	for k, v := range msg.EncryptedData.Nonces {
		new.EncryptedData.Nonces[k] = v
	}
	return new
}
