ARG GO_VERSION=1.20

FROM golang:${GO_VERSION}-alpine as build

RUN apk update
RUN apk add --no-cache build-base git

WORKDIR /stegosaurus

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o stegosaurus cmd/stegosaurus/main.go

FROM alpine

COPY --from=build stegosaurus .
COPY --from=build pipeline/encode .