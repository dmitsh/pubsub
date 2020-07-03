all: publisher server subscriber

server: cmd/server/main.go
	go build ./cmd/server

subscriber: cmd/subscriber/main.go
	go build ./cmd/subscriber

publisher: cmd/publisher/main.go
	go build ./cmd/publisher

clean:
	rm -f server subscriber publisher

.PHONY: clean
