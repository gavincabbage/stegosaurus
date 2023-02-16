ARG GO_VERSION=1.19

FROM golang:${GO_VERSION}-alpine

RUN apk update
RUN apk add --no-cache build-base git

WORKDIR /stegosaurus

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .