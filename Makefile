APP_NAME=app
GRPC_PKG=./cmd/app

.PHONY: build test docker-build run lint gofmt govet

build:
	go build -o $(APP_NAME) $(GRPC_PKG)

# Runs formatting + vet + tests
test: gofmt govet
	go test ./...

gofmt:
	gofmt -w .

govet:
	go vet ./...

lint:
	golangci-lint run ./...

docker-build:
	docker build -t usdt-rates:latest .

run:
	./$(APP_NAME)