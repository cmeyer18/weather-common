FROM golang:1.19-alpine

WORKDIR /app

COPY . .

RUN go mod vendor

RUN go build ./...
