# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate templ files and build
RUN go install github.com/a-h/templ/cmd/templ@latest && \
    templ generate && \
    go build -o server ./cmd/server

# Production stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

CMD ["./server"]