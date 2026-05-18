.PHONY: build test test-race lint tidy vet clean

build:
	go build -o bin/fxlint ./cmd/fxlint

test:
	go test ./...

test-race:
	go test -race ./...

vet:
	go vet ./...

lint:
	golangci-lint run

tidy:
	go mod tidy

clean:
	rm -rf bin/ coverage.*
