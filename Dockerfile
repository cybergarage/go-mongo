FROM alpine:latest
RUN apk update && apk add git go

USER root

COPY . /go-mongo
WORKDIR /go-mongo

RUN go build -o /go-mongod github.com/cybergarage/go-mongo/examples/go-mongod

ENTRYPOINT ["/go-mongod"]
