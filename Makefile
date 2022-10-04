.DEFAULT_GOAL := build

fmt:
	go fmt ./...

lint: fmt
	golint ./...

vet: fmt
	go vet ./...

verify: fmt
	go mod verify

build: vet verify
	go build -o ./build/pg-cache ./cmd

build-linux-amd64: vet verify
	GOOS=linux GOARCH=amd64 go build -o ./build/pg-cache-linux-amd64 ./cmd

test: vet verify
	go run ./cmd/main.go
