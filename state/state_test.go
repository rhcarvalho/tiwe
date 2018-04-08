package state

import (
	"fmt"
	"reflect"
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

func TestXor(t *testing.T) {
	tests := []struct {
		in   [][8]byte
		want []byte
	}{
		{
			in:   nil,
			want: nil,
		},
		{
			in:   [][8]byte{{1, 1, 1, 1}},
			want: []byte{1, 1, 1, 1, 0, 0, 0, 0},
		},
		{
			in:   [][8]byte{{1, 1, 1, 1}, {1, 1, 0, 1}},
			want: []byte{0, 0, 1, 0, 0, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%x", tt.in), func(t *testing.T) {
			if got := xor(tt.in...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got % x, want % x", got, tt.want)
			}
		})
	}
}

func TestXorAlloc(t *testing.T) {
	in := [][8]byte{{1, 1, 1, 1}, {1, 1, 0, 1}}
	want := [][8]byte{{1, 1, 1, 1}, {1, 1, 0, 1}}
	xor(in...)
	if !reflect.DeepEqual(in, want) {
		t.Errorf("input mutated: got %x, want %x", in, want)
	}
}
