# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.16-buster AS build

WORKDIR /volleybot

COPY * ./
RUN go mod download
RUN go build -o ./vbot

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build ./vbot ./vbot

USER nonroot:nonroot

ENTRYPOINT ["./vbot"]
