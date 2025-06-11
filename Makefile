build:
	go build -o bin/goland-distributed-cache

run: build
	./bin/goland-distributed-cache