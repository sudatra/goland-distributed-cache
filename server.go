package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

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
	var (
		rawStr = string(rawCmd)
		parts = strings.Split(rawStr, " ")
	)
	if len(parts) == 0 {
		log.Println("Invalid Command");
		return;
	}

	cmd := Command(parts[0]);
	if cmd == CMDSet {
		if len(parts) != 4 {
			log.Println("Invalid SET Command");
			return;
		}

		ttl, err := strconv.Atoi(parts[3]);
		if err != nil {
			log.Println("Invalid SET Command");
			return;
		}

		msg := MSGSet{
			Key: []byte(parts[1]),
			Value: []byte(parts[2]),
			TTL: time.Duration(ttl),
		}
		if err := s.handleSetCmd(conn, msg); err != nil {
			return;
		}
	}
}

func (s *Server) handleSetCmd(conn net.Conn, msg MSGSet) error {
	fmt.Println("Handling the set command: ", msg);
	return nil;
}