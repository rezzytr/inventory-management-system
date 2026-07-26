# Stage 1: Build Phase
FROM golang:alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary Go secara terkompilasi statis
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/api

# Stage 2: Minimal Runtime Image
FROM alpine:latest

WORKDIR /app

# Install ca-certificates & tzdata
RUN apk --no-cache add ca-certificates tzdata

# Copy binary & assets dari stage builder
COPY --from=builder /app/main .
COPY --from=builder /app/web ./web

EXPOSE 8080

CMD ["./main"]