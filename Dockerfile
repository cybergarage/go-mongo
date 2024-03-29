FROM golang:1.20-alpine

USER root

COPY . /go-mongo
WORKDIR /go-mongo

RUN go build -o /go-mongod github.com/cybergarage/go-mongo/examples/go-mongod

ENTRYPOINT ["/go-mongod"]
