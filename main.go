package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/sudatra/goland-distributed-cache/cache"
)

func main() {
	listenAddr := flag.String("listenaddr", ":3000", "listen address of the server ");
	leaderAddr := flag.String("leaderaddr", "", "listen address of the leader");
	flag.Parse();

	opts := ServerOpts{
		ListenAddr: *listenAddr,
		IsLeader:   true,
		LeaderAddr: *leaderAddr,
	}

	go func() {
		time.Sleep(time.Second * 2);

		conn, err := net.Dial("tcp", ":3000");
		if err != nil {
			log.Fatal(err);
		}

		conn.Write([]byte("SET Foo Bar 2500000000"));
		time.Sleep(time.Second * 2);
		conn.Write([]byte("GET Foo"));

		buff := make([]byte, 1000);
		n, _ := conn.Read(buff);
		fmt.Println(string(buff[:n]));
	}()

	server := NewServer(opts, cache.New());
  server.Start();
}