package main

type MessageType uint8

const (
	NewGame MessageType = iota
)

type Message struct {
	Type MessageType
}
