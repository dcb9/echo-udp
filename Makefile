.PHONY: all
all:
	go build -o server.bin ./server/
	go build -o client.bin ./client/
