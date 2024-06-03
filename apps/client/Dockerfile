# syntax=docker/dockerfile:1

FROM golang:1.22.2

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /client ./apps/client

ENTRYPOINT [ "/client" ]
