# Start from the base Golang image
FROM golang:1.19.4

# Set the current working directory inside the container
WORKDIR /go/src/github.com/harlow/go-micro-services

# Download dependencies first for layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Build the application
RUN go install -ldflags="-s -w" ./cmd/...
