# Multi-stage build for ABíbliaDigital Go API & Landing Page
FROM golang:alpine AS builder

WORKDIR /app

# Install git, ca-certificates, python3 and sqlite
RUN apk add --no-cache git ca-certificates python3 sqlite

# Copy go.mod and go.sum for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code, database scripts and static files
COPY . .

# Generate biblia.db SQLite database from JSON sources
RUN python3 scripts/convert_all_versions.py

# Build stateless production binary (Pure Go, static)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o server ./cmd/server/main.go

# Production runtime stage
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/server /app/server
COPY --from=builder /app/biblia.db /app/biblia.db
COPY --from=builder /app/public /app/public
COPY --from=builder /app/docs /app/docs

ENV PORT=8080
ENV DB_PATH=/app/biblia.db

EXPOSE 8080

ENTRYPOINT ["/app/server"]
