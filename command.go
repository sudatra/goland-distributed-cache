package main

import (
	"time"
)

type Command string

type Message struct {
	Cmd   Command
	Key   []byte
	Value []byte
	TTL   time.Duration
}

type MSGSet struct {
	Key   []byte
	Value []byte
	TTL   time.Duration
}

type MSGGet struct {
	Key   []byte
}

const (
	CMDSet Command = "SET"
	CMDGet Command = "GET"
)
