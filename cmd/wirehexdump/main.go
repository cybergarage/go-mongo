// Copyright (C) 2019 The go-mongo Authors. All rights reserved.
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
monghexdump is a dump utility for MongoDB packet hexdump file.

	NAME
	 monghexdump

	SYNOPSIS
	 monghexdump <BSON File>

	RETURN VALUE
	  Return EXIT_SUCCESS or EXIT_FAILURE
*/
//nolint:forbidigo
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cybergarage/go-mongo/mongo/bson"
	"github.com/cybergarage/go-mongo/mongo/protocol"
)

const (
	ProgramName = "monghexdump"
)

func usages() {
	println("Usage:")
	println("  " + ProgramName + " FILE")
	println("")
	println("Return Value:")
	println("  Return EXIT_SUCCESS or EXIT_FAILURE")
	os.Exit(1)
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) <= 0 {
		usages()
	}

	bsonFilename := args[0]

	protocolBytes, err := os.ReadFile(bsonFilename)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	msg, err := protocol.NewMessageWithBytes(protocolBytes)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	fmt.Printf("Header : %d %d %d %s(%d)\n",
		msg.MessageLength(),
		msg.RequestID(),
		msg.ResponseTo(),
		msg.OpCode().String(),
		msg.OpCode(),
	)

	for _, doc := range msg.Documents() {
		jsonStr, err := bson.DocumentToJSONString(doc)
		if err != nil {
			println(err.Error())
			os.Exit(1)
		}
		fmt.Println(jsonStr)
	}

	os.Exit(0)
}
