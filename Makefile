# Copyright (C) 2019 The go-mongo Authors. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http:#www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

SHELL := bash

PREFIX?=$(shell pwd)

GIT_ROOT=github.com/cybergarage/
MODULE_NAME=go-mongo

PKG_NAME=mongo
PKG_VER=$(shell git describe --abbrev=0 --tags)

MODULE_ROOT=${PKG_NAME}
MODULE_SRC_DIR=${PKG_NAME}
MODULE_SRCS=\
	${MODULE_SRC_DIR}/bson \
	${MODULE_SRC_DIR}/message \
	${MODULE_SRC_DIR}/protocol \
	${MODULE_SRC_DIR}/shell \
	${MODULE_SRC_DIR}
MODULE_PKG_ROOT=${GIT_ROOT}${MODULE_NAME}/${MODULE_ROOT}
MODULE_PKGS=\
	${MODULE_PKG_ROOT}/bson \
	${MODULE_PKG_ROOT}/message \
	${MODULE_PKG_ROOT}/protocol \
	${MODULE_PKG_ROOT}/shell \
	${MODULE_PKG_ROOT}

EXAMPLES_ROOT=examples
EXAMPLES_PKG_ROOT=${GIT_ROOT}${MODULE_NAME}/${EXAMPLES_ROOT}
EXAMPLES_DEAMON_BIN=go-mongod
EXAMPLES_DOCKER_TAG=cybergarage/${EXAMPLES_DEAMON_BIN}:${PKG_VER}
EXAMPLES_DEAMON_ROOT=${EXAMPLES_PKG_ROOT}/${EXAMPLES_DEAMON_BIN}
EXAMPLES_SRC_DIR=${EXAMPLES_ROOT}/${EXAMPLES_DEAMON_BIN}
EXAMPLES_SRCS=\
	${EXAMPLES_SRC_DIR} \
	${EXAMPLES_SRC_DIR}/server
EXAMPLES_PKGS=\
	${EXAMPLES_DEAMON_ROOT} \
	${EXAMPLES_DEAMON_ROOT}/server
EXAMPLE_BINARIES=\
	${EXAMPLES_DEAMON_ROOT}

BIN_ROOT=cmd
BIN_PKG_ROOT=${GIT_ROOT}${MODULE_NAME}/${BIN_ROOT}
BIN_SRC_DIR=${BIN_ROOT}
BIN_BINARIES=\
	${BIN_PKG_ROOT}/wirehexdump

TEST_ROOT=${PKG_NAME}test
TEST_PKG_NAME=${TEST_ROOT}
TEST_PKG_ROOT=${GIT_ROOT}${MODULE_NAME}/${TEST_PKG_NAME}
TEST_PKG_DIR=${TEST_PKG_NAME}
TEST_PKG_SRCS=\
        ${TEST_PKG_DIR}
TEST_PKGS=\
        ${TEST_PKG_ROOT}
TEST_FILE_DIR=${TEST_ROOT}/test
TEST_HELPER_NAME=embed
TEST_HELPER=${TEST_FILE_DIR}/${TEST_HELPER_NAME}.go

ALL_SRCS=\
	${MODULE_SRCS} \
	${TEST_PKG_SRCS} \
	${EXAMPLES_SRCS}

ALL_ROOTS=\
	${MODULE_ROOT} \
	${TEST_ROOT} \
	${EXAMPLES_ROOT} \
	${BIN_ROOT}

ALL_PKGS=\
	${MODULE_PKGS} \
	${TEST_PKGS} \
	${EXAMPLES_PKGS}

BINARIES=\
	${EXAMPLE_BINARIES} \
	${BIN_BINARIES}

.PHONY: version clean format vet lint
.IGNORE: test

all: test

version:
	@pushd ${MODULE_SRC_DIR} && ./version.gen > version.go && popd

${TEST_HELPER} : ${TEST_FILE_DIR}/${TEST_HELPER_NAME}.pl $(wildcard ${TEST_FILE_DIR}/*.qst)
	perl $< > $@

format: ${TEST_HELPER} 
	gofmt -s -w ${ALL_ROOTS}

vet: format
	go vet ${ALL_PKGS}

lint: vet
	golangci-lint run ${ALL_SRCS}

test: test
	go test -v -cover -p=1 ${ALL_PKGS}

install:
	go install -v -gcflags=${GCFLAGS} ${BINARIES}

build: test
	go build  -v -gcflags=${GCFLAGS} ${BINARIES}

run: build
	./${EXAMPLES_DEAMON_BIN}

image: test
	docker image build -t ${EXAMPLES_DOCKER_TAG} .

rund: image
	docker container run -it --rm -p 27017:27017 ${EXAMPLES_DOCKER_TAG}

mongod:
	docker container run -it --rm -p 27017:27017 mongo:4.4.19

clean:
	go clean -i ${ALL_PKGS}
