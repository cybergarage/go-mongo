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
PRODUCT_NAME=go-mongo
PACKAGE_NAME=mongo

MODULE_ROOT=${PACKAGE_NAME}
MODULE_SRC_DIR=${PACKAGE_NAME}
MODULE_SRCS=\
	${MODULE_SRC_DIR}/bson \
	${MODULE_SRC_DIR}/log \
	${MODULE_SRC_DIR}/message \
	${MODULE_SRC_DIR}/protocol \
	${MODULE_SRC_DIR}
MODULE_PACKAGE_ROOT=${GIT_ROOT}${PRODUCT_NAME}/${MODULE_ROOT}
MODULE_PACKAGES=\
	${MODULE_PACKAGE_ROOT}/bson \
	${MODULE_PACKAGE_ROOT}/log \
	${MODULE_PACKAGE_ROOT}/message \
	${MODULE_PACKAGE_ROOT}/protocol \
	${MODULE_PACKAGE_ROOT}

EXAMPLES_ROOT=examples
EXAMPLES_PACKAGE_ROOT=${GIT_ROOT}${PRODUCT_NAME}/${EXAMPLES_ROOT}
EXAMPLES_DEAMON_BIN=go-mongod
EXAMPLES_DEAMON_ROOT=${EXAMPLES_PACKAGE_ROOT}/${EXAMPLES_DEAMON_BIN}
EXAMPLES_SRC_DIR=${EXAMPLES_ROOT}/${EXAMPLES_DEAMON_BIN}
EXAMPLES_SRCS=\
	${EXAMPLES_SRC_DIR} \
	${EXAMPLES_SRC_DIR}/server
EXAMPLES_PACKAGES=\
	${EXAMPLES_DEAMON_ROOT} \
	${EXAMPLES_DEAMON_ROOT}/server
EXAMPLE_BINARIES=\
	${EXAMPLES_DEAMON_ROOT}

ALL_ROOTS=\
	${MODULE_ROOT} \
	${EXAMPLES_ROOT}

ALL_PACKAGES=\
	${MODULE_PACKAGES} \
	${EXAMPLES_PACKAGES}

.PHONY: clean format vet lint

all: test

format:
	gofmt -w ${ALL_ROOTS}

vet: format
	go vet ${ALL_PACKAGES}

lint: vet
	golangci-lint run ${MODULE_SRCS} ${EXAMPLES_SRCS}

build: lint
	go build -v ${MODULE_PACKAGES}

test:
	go test -v -cover -p=1 ${ALL_PACKAGES}

clean:
	go clean -i ${ALL_PACKAGES}