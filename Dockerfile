FROM golang:1.24.1-alpine3.21 AS base

WORKDIR /todo-api

COPY go.mod go.sum ./

RUN go mod download

COPY . .
CMD ["go","run","main.go"]