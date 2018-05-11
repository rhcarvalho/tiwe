package main

import "errors"

// PTH stands for Pool, Table and Hands.
type PTH struct {
	Pool
	Table
	Hands map[string][]Tile
}

// Transition represents a change in a PTH.
type Transition struct {
	PreviousHash string
	Pool
	Table
	Player string
	Hand   []Tile
}

// A Pool is a collection of tiles whose faces are unknown to all players.
type Pool []ConcealedTile

// A Table is a collection of sets (runs or groups). Tiles on the Table are
// known to all players.
type Table []Set

type Set interface{}

type Run []Tile

type Group []Tile

// A ConcealedTile represents a Tile whose face value is concealed to one or
// more players.
type ConcealedTile struct {
}

// Seal encrypts the tile with one or more secret keys. It is okay to call Seal
// multiple times. Sealing with the same key more than once has the same effect
// as sealing once.
func (ct *ConcealedTile) Seal(keys ...[]byte) {

}

// Open decrypts the tile with one of more secret keys. It is okay to call Open
// multiple times. Calling Open with a key that is innefective is a no-op.
func (ct *ConcealedTile) Open(keys ...[]byte) {

}

// Tile returns a Tile if and only if the ConcealedTile can be revealed. To
// reveal a ConcealedTile, call Open with the same keys used to seal it, in any
// order.
func (ct *ConcealedTile) Tile() (*Tile, error) {
	return nil, errors.New("could not reveal tile")
}

type Tile struct {
	Value uint64
	Color uint64
}
