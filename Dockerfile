# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/a-h/templ/cmd/templ@latest && \
    templ generate && \
    CGO_ENABLED=0 GOOS=linux go build -o /sidekick cmd/server/main.go cmd/server/routes.go

# Final stage
FROM alpine:3.18

WORKDIR /app

COPY --from=builder /sidekick .
COPY --from=builder /app/web/templates ./web/templates
COPY --from=builder /app/web/static ./web/static

EXPOSE 3000

CMD ["./sidekick"]
