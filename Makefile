.PHONY: lint start build

lint:
	golangci-lint run
start:
	go run cmd/hexlet-path-size/main.go		
build:
	go build -o bin/hexlet-path-size ./cmd/hexlet-path-size	
	bin/hexlet-path-size