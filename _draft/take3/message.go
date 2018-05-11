package main

import "time"

type Log []Message

type Message struct {
	ParentHash string
	ID         string
	From       string
	Type       MessageType
	Bytes      []byte
}

type MessageType uint8

const (
	MessageTypeUndefined MessageType = iota
	MessageTypeRequestToPlay
	MessageTypeAcceptToPlay
	MessageTypeUpdatePTH
	MessageTypeRevealSecret
)

type RequestToPlay struct {
	ID       string
	Deadline time.Time
}
