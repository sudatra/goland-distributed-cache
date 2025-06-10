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
}