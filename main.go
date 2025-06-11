package main

import (
	"log"
	"net"
	"time"

	"github.com/sudatra/goland-distributed-cache/cache"
)

func main() {
	opts := ServerOpts{
		listenAddr: ":3000",
		isLeader:   true,
	}

	go func() {
		time.Sleep(time.Second * 2);

		conn, err := net.Dial("tcp", ":3000");
		if err != nil {
			log.Fatal(err);
		}

		conn.Write([]byte("SET Foo Bar 2500"));
	}()

	server := NewServer(opts, cache.New());
  server.Start();
}