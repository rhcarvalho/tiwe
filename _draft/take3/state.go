package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"errors"
)

type state struct {
	Players []string
	Log
	PTH
	Secrets map[string][][]byte // hash of encrypted title -> []secret
}

func (s *state) HandleMessage(msg Message) error {
	// TODO: transactional changes to state: either update Log and PTH or
	// nothing.
	err := s.AppendMessage(msg)
	if err != nil {
		return err
	}
	switch msg.Type {
	case MessageTypeUpdatePTH:
		dec := gob.NewDecoder(bytes.NewReader(msg.Bytes))
		var t Transition
		err := dec.Decode(&t)
		if err != nil {
			return err
		}
		return s.UpdatePTH(t)
	}
	return nil
}

func (s *state) AppendMessage(msg Message) error {
	if len(s.Log) == 0 {
		s.Log = append(s.Log, msg)
		return nil
	}
	lastMsg := s.Log[len(s.Log)-1]
	h, err := marshalHash(lastMsg)
	if err != nil {
		return err
	}
	if msg.ParentHash != h {
		return errors.New("bad parent hash")
	}
	s.Log = append(s.Log, msg)
	return nil
}

func (s *state) UpdatePTH(t Transition) error {
	h, err := marshalHash(s.PTH)
	if err != nil {
		return err
	}
	if t.PreviousHash != h {
		return errors.New("bad transition")
	}
	// TODO: validate new PTH
	s.PTH.Pool = t.Pool
	s.PTH.Table = t.Table
	// TODO: validate that it is t.Player's turn
	s.PTH.Hands[t.Player] = t.Hand
	return nil
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

func hash(b []byte) string {
	sum := sha256.Sum256(b)
	return string(sum[:])
}

func marshalHash(v interface{}) (string, error) {
	b, err := marshal(v)
	if err != nil {
		return "", err
	}
	return hash(b), nil
}

func test() {
	log := Log{
		// Assume:
		// 1. Implicit initial play order is Alice, Bob, Carol
		// 2. They can communicate with each other

		// Phase: define play order
		Message{
			ParentHash: "",
			ID:         "1",
			From:       "Alice",
			Type:       MessageTypeUpdatePTH,
			Bytes: mustMarshal(Transition{
				PreviousHash: "",
				Pool:         Pool{}, // TODO: shuffle-encrypted tiles 1..N
			}),
		},
		Message{
			ParentHash: "...",
			ID:         "2",
			From:       "Bob",
			Type:       MessageTypeUpdatePTH,
			Bytes: mustMarshal(Transition{
				PreviousHash: "...",
				Pool:         Pool{}, // TODO: shuffle-encrypted tiles 1..N
			}),
		},
		Message{
			ParentHash: "...",
			ID:         "3",
			From:       "Carol",
			Type:       MessageTypeUpdatePTH,
			Bytes: mustMarshal(Transition{
				PreviousHash: "...",
				Pool:         Pool{}, // TODO: shuffle-encrypted tiles 1..N
			}),
		},
		// Subphase: reveal secrets
		Message{
			ParentHash: "...",
			ID:         "4",
			From:       "Alice",
			Type:       MessageTypeRevealSecret,
			Bytes:      mustMarshal([]byte{ /* TODO: Alice's secret */ }),
		},
		Message{
			ParentHash: "...",
			ID:         "5",
			From:       "Bob",
			Type:       MessageTypeRevealSecret,
			Bytes:      mustMarshal([]byte{ /* TODO: Bob's secret */ }),
		},
		Message{
			ParentHash: "...",
			ID:         "6",
			From:       "Carol",
			Type:       MessageTypeRevealSecret,
			Bytes:      mustMarshal([]byte{ /* TODO: Carol's secret */ }),
		},
		// At this point, all players can decode Pool and figure out the
		// new play order.

		// Assume new order is Bob, Alice, Carol

		// Phase: initial shuffle
		Message{
			ParentHash: "...",
			ID:         "###",
			From:       "Bob",
			Type:       MessageTypeUpdatePTH,
			Bytes: mustMarshal(Transition{
				PreviousHash: "...",
				Pool:         Pool{}, // TODO: shuffle-encrypted tiles 1..106
			}),
		},
		Message{
			ParentHash: "...",
			ID:         "###",
			From:       "Alice",
			Type:       MessageTypeUpdatePTH,
			Bytes: mustMarshal(Transition{
				PreviousHash: "...",
				Pool:         Pool{}, // TODO: shuffle-encrypted tiles 1..106
			}),
		},
		Message{
			ParentHash: "...",
			ID:         "###",
			From:       "Carol",
			Type:       MessageTypeUpdatePTH,
			Bytes: mustMarshal(Transition{
				PreviousHash: "...",
				Pool:         Pool{}, // TODO: shuffle-encrypted tiles 1..106
			}),
		},
		// Now the Pool is shuffled, nobody knows where any particular
		// tile is.

		// Phase: each player takes 14 tiles
		Message{
			ParentHash: "...",
			ID:         "###",
			From:       "Bob",
			Type:       MessageTypeUpdatePTH,
			Bytes: mustMarshal(Transition{
				PreviousHash: "...",
				Pool:         Pool{}, // TODO: shuffle-encrypted tiles 1..106 - 14
				Player:       "Bob",
				Hand:         []Tile{ /* TODO: 14 tiles chosen by Bob */ },
			}),
		},
		Message{
			ParentHash: "...",
			ID:         "###",
			From:       "Alice",
			Type:       MessageTypeUpdatePTH,
			Bytes: mustMarshal(Transition{
				PreviousHash: "...",
				Pool:         Pool{}, // TODO: shuffle-encrypted tiles 1..106 -14 -14
				Player:       "Alice",
				Hand:         []Tile{ /* TODO: 14 tiles chosen by Alice */ },
			}),
		},
		Message{
			ParentHash: "...",
			ID:         "###",
			From:       "Carol",
			Type:       MessageTypeUpdatePTH,
			Bytes: mustMarshal(Transition{
				PreviousHash: "...",
				Pool:         Pool{}, // TODO: shuffle-encrypted tiles 1..106 -14 -14 -14
				Player:       "Carol",
				Hand:         []Tile{ /* TODO: 14 tiles chosen by Carol */ },
			}),
		},
		// Subphase: reveal secrets
		// Bob reveals 28 secrets (14 -> Alice's hand, 14 -> Carol's hand)
		// Alice reveals 28 secrets (...)
		// Carol reveals 28 secrets (...)

		// Now all players have an initial hand. They can either play
		// 50+ points into the table as their initial meld, or take 1
		// tile from the Pool.

		// Suppose Bob doesn't have an initial meld to play, he goes
		// take 1 tile from the Pool.
		Message{
			ParentHash: "...",
			ID:         "###",
			From:       "Bob",
			Type:       MessageTypeUpdatePTH,
			Bytes: mustMarshal(Transition{
				PreviousHash: "...",
				Pool:         Pool{}, // TODO: previous Pool -1 tile
				Player:       "Bob",
				Hand:         []Tile{ /* TODO: 14+1 tiles chosen by Bob */ },
			}),
		},
		// Bob doesn't know what tile he took from the Pool until Alice
		// and Carol reveal secrets.
		// They must do that in their next turn!
	}
	_ = log
}
