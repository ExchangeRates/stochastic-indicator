.PHONY: build
build:
	go fmt ./...
	go build -v ./cmd/server

.PHONY: test
test:
	go test -v -race -time 30s ./...

.DEFAULT_GOAL := build
