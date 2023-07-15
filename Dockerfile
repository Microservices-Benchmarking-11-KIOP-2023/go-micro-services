FROM golang:1.20.6

WORKDIR /go/src/github.com/harlow/go-micro-services

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install -ldflags="-s -w" ./cmd/...
