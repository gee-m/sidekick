.PHONY: dev build test clean docker-prod docker-dev

# Development
dev:
	docker-compose up --build

# Production
docker-prod:
	docker build -t sidekick:latest -f Dockerfile .

# Testing
test:
	go test -v ./...

# Build
build:
	go build -o bin/sidekick cmd/server/main.go

# Clean
clean:
	rm -rf bin/
	docker-compose down -v

# Generate templ files
generate:
	templ generate

# Database
db-migrate:
	go run cmd/migrate/main.go up

db-rollback:
	go run cmd/migrate/main.go down

# Air (local development without Docker)
air:
	air -c .air.toml

# Development dependencies
dev-deps:
	go install github.com/air-verse/air@latest
	go install github.com/a-h/templ/cmd/templ@latest
