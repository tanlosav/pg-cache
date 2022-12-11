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


build-schema: check
	go build -o ./build/pg-schema ./cmd/schema

build-schema-linux-amd64: check
	GOOS=linux GOARCH=amd64 go build -o ./build/pg-schema-linux-amd64 ./cmd/schema

test-schema: check
	go run ./cmd/schema/main.go --configuration-provider file --configuration-source ./test/config.yml


build-cache: check
	go build -o ./build/pg-cache ./cmd/cache

build-cache-linux-amd64: check
	GOOS=linux GOARCH=amd64 go build -o ./build/pg-cache-linux-amd64 ./cmd/cache

test-cache: check
	go run ./cmd/cache/main.go --configuration-provider file --configuration-source ./test/config.yml
