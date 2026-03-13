.PHONY: dev build test sqlc clean

# Development
dev:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server
	go build -o bin/seed ./cmd/seed

# Testing
test:
	go test ./internal/domain/... -v

# Code generation
sqlc:
	sqlc generate

# User management
seed-create:
	go run ./cmd/seed create --email $(EMAIL) --password $(PASSWORD) --name "$(NAME)"

seed-list:
	go run ./cmd/seed list

# Cleanup
clean:
	rm -f bin/server bin/seed

# Estimate Builder (Svelte)
ui-install:
	cd ui && npm install

ui:
	cd ui && npm run build

ui-watch:
	cd ui && npm run watch

# Dependencies
deps:
	go mod download
	go mod tidy
