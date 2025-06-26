#!/bin/bash
# Build stage
FROM golang:1.19-alpine3.16 AS builder
WORKDIR /app
RUN apk update && apk add --no-cache alpine-sdk
COPY . .
# RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -o wms-server -tags musl

# Run stage
FROM alpine:3.16
WORKDIR /app

RUN apk update && apk add --no-cache git
RUN apk update && apk add --no-cache tzdata
RUN apk update && apk add --no-cache gcompat
ENV TZ="Asia/Jakarta"

COPY --from=builder /app/wms-server .

EXPOSE 7878
ENTRYPOINT [ "/app/wms-server" ]
# ENTRYPOINT [ "/app/start.sh" ]