package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/sudatra/goland-distributed-cache/cache"
)

type ServerOpts struct {
	ListenAddr string
	IsLeader   bool
	LeaderAddr string
}

type Server struct {
	ServerOpts
	followers map[net.Conn]struct{}
	cache			cache.Cacher
}

func NewServer(opts ServerOpts, c cache.Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		cache: c,
		// TODO: only allocate when we are the leader
		followers: make(map[net.Conn]struct{}),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr);
	if err != nil {
		return fmt.Errorf("Listen error: (%s)", err);
	}

	log.Printf("Server starting on port [%s]\n ", s.ListenAddr);

	if !s.IsLeader {
		go func ()  {
			conn, err := net.Dial("tcp", s.LeaderAddr);
			fmt.Println("Connected with leader: ", s.LeaderAddr);
			if err != nil {
				log.Fatal(err);
			}

			s.handleConn(conn);
		}()
	}

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

	if s.IsLeader {
		s.followers[conn] = struct{}{};
	}
	fmt.Println("connection made: ", conn.RemoteAddr());

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
		conn.Write([]byte(err.Error()));
		return;
	}

	fmt.Printf("received command %s", msg.Cmd);
	
	switch msg.Cmd {
		case CMDSet:
			err = s.handleSetCmd(conn, msg);
		case CMDGet:
			err = s.handleGetCmd(conn, msg);
	}

	if err != nil {
		fmt.Println("Failed to handle command", err);
		conn.Write([]byte(err.Error()));
	}
}

func (s *Server) handleSetCmd(conn net.Conn, msg *Message) error {
	if err := s.cache.Set(msg.Key, msg.Value, msg.TTL); err != nil {
		return err;
	}
	go s.sendToFollowers(context.TODO(), msg);

	return nil;
}

func (s *Server) handleGetCmd(conn net.Conn, msg *Message) error {
	val, err := s.cache.Get(msg.Key);
	if err != nil {
		return err;
	}

	_, err = conn.Write(val);
	return err;
}

func (s *Server) sendToFollowers(ctx context.Context, msg *Message) error {
	for conn := range s.followers {
		fmt.Println("Forwading key to follower");
		rawMsg := msg.ToBytes();
		fmt.Println("Forwading rawMsg to follower: ", string(rawMsg));

		_, err := conn.Write(rawMsg);
		if err != nil {
			fmt.Println("write to follower error ", err);
			continue;
		}
	}
	return nil;
}