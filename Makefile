.PHONY: lint start build

test:
	 go test ./... 
lint:
	@echo "Running golangci-lint in container..."
	@docker compose exec -e GOFLAGS="-buildvcs=false" golangci-lint run --timeout=5m ./...

start:
	go run cmd/hexlet-path-size/main.go		
build:
	go build -o bin/hexlet-path-size ./cmd/hexlet-path-size	
	bin/hexlet-path-size