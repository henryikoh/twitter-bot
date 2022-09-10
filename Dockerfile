# syntax=docker/dockerfile:1

FROM golang:alpine3.16

RUN apk add build-base

WORKDIR /app

COPY . .

RUN go mod download

WORKDIR /app

RUN go build -o /twitter-bot


EXPOSE 5000

CMD [ "/twitter-bot" ]