FROM golang:1.24-alpine

WORKDIR /app

RUN apk add --no-cache file && \
    echo "🔍 Binary format:" && \
    file ./MaliciousBlocker

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o MaliciousBlocker ./MaliciousBlocker.go
RUN chmod +x ./MaliciousBlocker

# 🔍 Inspect the binary format
RUN file ./MaliciousBlocker

EXPOSE 9090

CMD ["./MaliciousBlocker"]