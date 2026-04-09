APP_NAME=app
GRPC_PKG=./cmd/app

.PHONY: build test docker-build run lint

build:
	go build -o $(APP_NAME) $(GRPC_PKG)

test:
	go test ./...

docker-build:
	docker build -t usdt-rates:latest .

run:
	./$(APP_NAME)

lint:
	golangci-lint run ./...
