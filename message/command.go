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

package message

import (
	"github.com/cybergarage/go-mongo/bson"
)

// See : Database Commands
// https://docs.mongodb.com/manual/reference/command/

const (
	//
	// Replication Commands
	//
	//isMaster displays information about this member’s role in the replica set, including whether it is the master.
	IsMaster = "isMaster"
)

// Command represents a query command of MongoDB database command.
type Command struct {
	bson.Element
}

// NewCommandWithElement returns a new command instance with the specified a BSON element.
func NewCommandWithElement(element bson.Element) *Command {
	cmd := &Command{
		Element: element,
	}
	return cmd
}

// CommandExecutor represents an interface for MongoDB database commands.
type CommandExecutor interface {
	// Replication Commands

	// ExecuteCommand handles query commands other than those explicitly specified above.
	ExecuteCommand(cmd *Command) (bson.Document, error)
}
