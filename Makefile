.DEFAULT_GOAL := build

fmt:
	go fmt ./...

vet:
	go vet ./...

staticcheck:
	staticcheck ./...

verify: fmt
	go mod verify

check: fmt vet staticcheck

build: check
	go build -o ./build/pg-cache ./cmd

build-linux-amd64: check
	GOOS=linux GOARCH=amd64 go build -o ./build/pg-cache-linux-amd64 ./cmd

test: check
	go run ./cmd/main.go
