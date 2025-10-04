.PHONY: build run clean test

build:
	go build -o bin/app cmd/main.go

run:
	go run cmd/main.go

clean:
	rm -rf bin/

test:
	go test ./...

deps:
	go mod download
	go mod tidy

