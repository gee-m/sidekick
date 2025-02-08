.PHONY: dev build test clean prod kill-8080 migrate-create migrate-up migrate-down

# Development
dev: kill-8080
	docker compose up --build

# Build development
build:
	docker compose build

# Build production
prod:
	docker build -t appgents:latest .

# Test
test:
	go test -v ./...

# Clean
clean:
	docker compose down -v
	rm -rf tmp/

# Generate templ files
generate:
	templ generate

# Install dependencies
deps:
	go install github.com/air-verse/air@latest
	go install github.com/a-h/templ/cmd/templ@latest

# Kill process on port 8080
kill-8080:
	@lsof -ti:8080 | xargs kill -9 2>/dev/null || true

# Create a new migration
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir internal/platform/database/migrations/files -seq $$name

# Run migrations up
migrate-up:
	migrate -path internal/platform/database/migrations/files -database "postgres://appgents:appgents@localhost:5432/appgents?sslmode=disable" up

# Run migrations down
migrate-down:
	migrate -path internal/platform/database/migrations/files -database "postgres://appgents:appgents@localhost:5432/appgents?sslmode=disable" down