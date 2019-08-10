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
	"fmt"
	"strings"

	"github.com/cybergarage/go-mongo/bson"
	"github.com/cybergarage/go-mongo/protocol"
)

// See : Database Commands
// https://docs.mongodb.com/manual/reference/command/

const (
	errorUnknownCommand = "Unknown Command : {%s}"
)

const (
	adminCommand = "admin.$cmd"
)

var allSupportedCommands = []string{
	IsMaster,
}

// Command represents a query command of MongoDB database command.
type Command struct {
	IsAdmin  bool
	Elements []bson.Element
}

// CommandExecutor represents an interface for MongoDB database commands.
type CommandExecutor interface {
	// ExecuteCommand handles query commands other than those explicitly specified above.
	ExecuteCommand(cmd *Command) ([]bson.Document, error)
}

// NewCommandWithQuery returns a new command instance with the specified BSON document.
func NewCommandWithQuery(q *protocol.Query) (*Command, error) {
	var err error
	var elements []bson.Element

	doc := q.GetQuery()
	if 0 <= len(doc) {
		elements, err = doc.Elements()
		if err != nil {
			return nil, err
		}
	}

	cmd := &Command{
		IsAdmin:  q.IsCollection(adminCommand),
		Elements: elements,
	}
	return cmd, nil
}

// IsAdminCommand returns true when it is a admin command, otherwise false.
func (cmd *Command) IsAdminCommand() bool {
	return cmd.IsAdmin
}

// GetType returns a string type
func (cmd *Command) GetType() (string, error) {
	for _, typeString := range allSupportedCommands {
		if cmd.IsType(typeString) {
			return typeString, nil
		}
	}
	return "", fmt.Errorf(errorUnknownCommand, cmd.String())
}

// IsType returns true when the command has the specified element, otherwise false.
func (cmd *Command) IsType(typeString string) bool {
	for _, element := range cmd.Elements {
		key := element.Key()
		if strings.ToLower(key) == typeString {
			return true
		}
	}
	return false
}

// String returns the string description.
func (cmd *Command) String() string {
	str := ""
	for n, element := range cmd.Elements {
		if n != 0 {
			str += " "
		}
		str += element.Key()
	}
	return str
}
