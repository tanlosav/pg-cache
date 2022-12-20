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

test: check
	go test ./...


#
# partitioner
#
build-partitioner: check
	go build -o ./build/partitioner ./cmd/partitioner

build-partitioner-linux-amd64: check
	GOOS=linux GOARCH=amd64 go build -o ./build/partitioner-linux-amd64 ./cmd/partitioner

run-partitioner: check
	go run ./cmd/partitioner/main.go --configuration-provider file --configuration-source ./test/config.yml


#
# rest-api-gateway
#
build-rest-api-gateway: check
	go build -o ./build/rest-api-gateway ./cmd/rest-api-gateway

build-rest-api-gateway-linux-amd64: check
	GOOS=linux GOARCH=amd64 go build -o ./build/rest-api-gateway-linux-amd64 ./cmd/rest-api-gateway

run-rest-api-gateway: check
	go run ./cmd/rest-api-gateway/main.go --configuration-provider file --configuration-source ./test/config.yml


#
# dev env
#
dev-env-start:
	./scripts/dev-env-start.sh

dev-env-ps:
	./scripts/dev-env-ps.sh

dev-env-down:
	./scripts/dev-env-down.sh
