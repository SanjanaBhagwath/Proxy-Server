FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum from this project folder
COPY go.mod .
RUN go mod download

# Copy the RateLimiter.go file
COPY RateLimiter.go .

# Build the Go application binary, specifying RateLimiter.go
RUN go build -o rate_limiter_proxy RateLimiter.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/rate_limiter_proxy .

EXPOSE 9090
ENTRYPOINT ["./rate_limiter_proxy"]