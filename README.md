# go-mongo

![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/cybergarage/go-mongo)
[![test](https://github.com/cybergarage/go-mongo/actions/workflows/make.yml/badge.svg)](https://github.com/cybergarage/go-mongo/actions/workflows/make.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/cybergarage/go-mongo.svg)](https://pkg.go.dev/github.com/cybergarage/go-mongo)
[![codecov](https://codecov.io/gh/cybergarage/go-mongo/graph/badge.svg?token=KTU8MELXYB)](https://codecov.io/gh/cybergarage/go-mongo)

The go-mongo is a database framework for implementing a [MongoDB](https://www.mongodb.com)-compatible server using Go easily.

## What is the go-mongo?

The go-mongo handles [MongoDB Wire Protocol](https://www.mongodb.com/docs/manual/reference/mongodb-wire-protocol/) and interprets the major messages automatically so that all developers can develop MongoDB-compatible servers easily. 

![](doc/img/framework.png)

Since the go-mongo handles all system commands automatically, developers can easily implement their MongoDB-compatible server only by simply handling user query commands.

## Table of Contents

- [Getting Started](doc/getting-started.md)

## Examples

- [Examples](doc/examples.md)
	- [go-mongod](examples/go-mongod) [![Docker Image Version](https://img.shields.io/docker/v/cybergarage/go-mongod)](https://hub.docker.com/repository/docker/cybergarage/go-mongod/)
	- [go-sqlserver](https://github.com/cybergarage/go-sqlserver) [![Docker Image Version](https://img.shields.io/docker/v/cybergarage/go-sqlserver)](https://hub.docker.com/repository/docker/cybergarage/go-sqlserver/)
	- [PuzzleDB](https://github.com/cybergarage/puzzledb-go) [![Docker Image Version](https://img.shields.io/docker/v/cybergarage/puzzledb)](https://hub.docker.com/repository/docker/cybergarage/puzzledb/)


## References

- [MongoDB](https://www.mongodb.com)
- [MongoDB Wire Protocol — MongoDB Manual](https://www.mongodb.com/docs/manual/reference/mongodb-wire-protocol/)
