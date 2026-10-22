package models

import (
	"net"
	"time"
)

type MessageType int

const (
	NewConnection MessageType = iota + 1
	NewMessage
	Delete
)

type Message struct {
	MessageType MessageType
	Conn        net.Conn
	Text        string
}

type Client struct {
	Conn        net.Conn
	LastMessage time.Time
	StrikCount  int
}
