build:
	go build -o bin/goland-distributed-cache

run: build
	./bin/goland-distributed-cache

runFollower: build
	./bin/goland-distributed-cache --listenaddr :4000 --leaderaddr :3000