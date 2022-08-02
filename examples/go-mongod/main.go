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

package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/cybergarage/go-mongo/examples/go-mongod/server"
	"github.com/cybergarage/go-mongo/mongo/log"
)

//////////////////////////////////////////////////
// main
//////////////////////////////////////////////////

func main() {
	log.SetStdoutDebugEnbled(true)

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
