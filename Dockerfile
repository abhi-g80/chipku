# syntax=docker/dockerfile:1
FROM golang:1.22-alpine

# Define build env
ENV GOOS linux
ENV CGO_ENABLED 0

WORKDIR /app

COPY . .
RUN go mod download

RUN apk add make gcc
RUN make test
RUN make build

EXPOSE 8090

ENTRYPOINT [ "build/chipku", "serve", "--port=8090" ]
