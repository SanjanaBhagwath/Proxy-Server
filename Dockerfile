# Use Golang base image
FROM golang:1.24

# Set the working directory inside the container
WORKDIR /app

# Copy go module files first
COPY go.mod go.sum ./
RUN go mod download

# Copy the Go source files
COPY . .

# Build the binary inside the container for the correct architecture
RUN go build -o ProxyServer ProxyServer.go

# Expose the required port
EXPOSE 9090

# Run the Proxy Server
CMD ["./ProxyServer"]