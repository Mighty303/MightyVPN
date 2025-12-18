.PHONY: build build-client build-server run-server run-client test clean

BINARY_DIR=bin
CLIENT_BINARY=$(BINARY_DIR)/client
SERVER_BINARY=$(BINARY_DIR)/server

all: build

build: build-client build-server

build-client:
	@mkdir -p $(BINARY_DIR)
	go build -o $(CLIENT_BINARY) ./cmd/client

build-server:
	@mkdir -p $(BINARY_DIR)
	go build -o $(SERVER_BINARY) ./cmd/server

run-server: build-server
	sudo $(SERVER_BINARY)

run-client: build-client
	sudo $(CLIENT_BINARY)

test:
	go test -v ./...

clean:
	rm -rf $(BINARY_DIR)

help:
	@echo "Available targets:"
	@echo "  make build         - Build both client and server"
	@echo "  make build-client  - Build client only"
	@echo "  make build-server  - Build server only"
	@echo "  make run-server    - Build and run server"
	@echo "  make run-client    - Build and run client"
	@echo "  make test          - Run tests"
	@echo "  make clean         - Clean build artifacts"