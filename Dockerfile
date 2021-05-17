FROM golang:alpine AS builder

ENV http_proxy="http://172.16.141.15:3128"
ENV APP_ENV="docker"
ENV GIN_MODE="release"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
COPY . .

ENV http_proxy=""

RUN go build -o main .

EXPOSE 7002

CMD ["./main"]

