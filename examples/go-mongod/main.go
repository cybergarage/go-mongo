// Copyright (C) 2022 The go-mongo Authors All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
go-mongod is an example of a compatible MongoDB server implementation using go-mongo.

	NAME
		go-mongod

	SYNOPSIS
		go-mongod [OPTIONS]

	OPTIONS
	-v      : Enable verbose output.
	-p      : Enable profiling.

	RETURN VALUE
		Return EXIT_SUCCESS or EXIT_FAILURE
*/

package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	clog "github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mysql/examples/go-mysqld/server"
)

//////////////////////////////////////////////////
// main
//////////////////////////////////////////////////

func main() {
	isDebugEnabled := flag.Bool("debug", false, "enable debugging log output")
	isProfileEnabled := flag.Bool("profile", false, "enable profiling server")
	flag.Parse()

	if *isDebugEnabled {
		clog.SetStdoutDebugEnbled(true)
	}

	if *isProfileEnabled {
		go func() {
			// nolint: gosec
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	server := server.NewServer()

	if err := server.Start(); err != nil {
		os.Exit(1)
	}

	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM)

	exitCh := make(chan int)

	go func() {
		for {
			s := <-sigCh
			switch s {
			case syscall.SIGHUP:
				if err := server.Restart(); err != nil {
					os.Exit(1)
				}
			case syscall.SIGINT, syscall.SIGTERM:
				if err := server.Stop(); err != nil {
					os.Exit(1)
				}
				exitCh <- 0
			}
		}
	}()

	code := <-exitCh

	os.Exit(code)
}
