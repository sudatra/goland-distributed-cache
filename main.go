package main

import (
	"github.com/sudatra/goland-distributed-cache/cache"
)

func main() {
	opts := ServerOpts{
		listenAddr: ":3000",
		isLeader:   true,
	}

	server := NewServer(opts, cache.New());

}