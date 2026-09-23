.PHONY: dev build test lint fmt

# Rebuild the site and restart the server on :8080 whenever files change.
dev:
	air

build:
	pnpm -C web build
	go build -o bin/gordle ./cmd/gordle

test:
	go test ./...

lint:
	go vet ./...
	pnpm -C web check

fmt:
	go fmt ./...
	pnpm -C web format
