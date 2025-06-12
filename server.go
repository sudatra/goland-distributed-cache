package main

import (
	"fmt"
	"log"
	"net"

	"github.com/sudatra/goland-distributed-cache/cache"
)

type ServerOpts struct {
	listenAddr string
	isLeader   bool
}

type Server struct {
	ServerOpts
	cache			cache.Cacher
}

func NewServer(opts ServerOpts, c cache.Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		cache: c,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr);
	if err != nil {
		return fmt.Errorf("Listen error: (%s)", err);
	}

	log.Printf("Server starting on port [%s]\n ", s.listenAddr);

	for {
		conn, err := ln.Accept();
		if err != nil {
			log.Printf("Accept error: %s\n", err);
			continue;
		}

		go s.handleConn(conn);
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer func() {
		conn.Close();
	}()

	buff := make([]byte, 2048);
	for {
		n, err := conn.Read(buff);
		if err != nil {
			log.Printf("Conn read error: %s\n", err);
			break;
		}

		msg := buff[:n];
		fmt.Printf(string(msg));

		go s.handleCommand(conn, buff[:n]);
	}
}

func (s *Server) handleCommand(conn net.Conn, rawCmd []byte) {
	msg, err := parseMessage(rawCmd);
	if err != nil {
		fmt.Println("Failed to parse command", err);
		return;
	}
	
	switch msg.Cmd {
		case CMDSet:
			if err := s.handleSetCmd(conn, msg); err != nil {
				return;
			}
	}
}

func (s *Server) handleSetCmd(conn net.Conn, msg *Message) error {
	fmt.Println("\nHandling the set command: ", msg);
	return nil;
}