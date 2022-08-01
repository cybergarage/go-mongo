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
	"github.com/cybergarage/go-mongo/mongo/bson"
	"github.com/cybergarage/go-mongo/mongo/protocol"
	"strings"
)

// See : Database Commands
// https://docs.mongodb.com/manual/reference/command/

const (
	errorUnknownCommand = "Unknown Command : {%s}"
)

const (
	adminCommand = "admin.$cmd"
)

const (
	IsMaster     = "ismaster"
	BuildInfo    = "buildinfo"
	GetLastError = "getlasterror"
)

// Command represents a query command of MongoDB database command.
type Command struct {
	IsAdmin  bool
	Elements []bson.Element
	Type     string
}

// CommandExecutor represents an interface for MongoDB database commands.
type CommandExecutor interface {
	// ExecuteCommand handles query commands other than those explicitly specified above.
	ExecuteCommand(cmd *Command) (bson.Document, error)
}

// NewCommandWithDocument returns a new command instance with the specified BSON document.
func NewCommandWithDocument(doc bson.Document) (*Command, error) {
	elements, err := doc.Elements()
	if err != nil {
		return nil, err
	}

	var cmdType string
	if 0 < len(elements) {
		// Example : "isMaster" or "ismaster"
		cmdType = strings.ToLower(elements[0].Key())
	}

	cmd := &Command{
		IsAdmin:  false,
		Elements: elements,
		Type:     cmdType,
	}
	return cmd, nil
}

// NewCommandWithQuery returns a new command instance with the specified BSON document.
func NewCommandWithQuery(q *protocol.Query) (*Command, error) {
	cmd, err := NewCommandWithDocument(q.GetQuery())
	if err != nil {
		return nil, err
	}

	cmd.IsAdmin = q.IsCollection(adminCommand)

	return cmd, nil
}

// NewCommandWithMsg returns a new command instance with the specified BSON document.
func NewCommandWithMsg(msg *protocol.Msg) (*Command, error) {
	cmd, err := NewCommandWithDocument(msg.GetBody())
	if err != nil {
		return nil, err
	}

	isAdmin := false

	bodyDoc := msg.GetBody()
	dbVal, err := bodyDoc.LookupErr("$db")
	if err == nil {
		dbStr, ok := dbVal.StringValueOK()
		if ok {
			if dbStr == "admin" {
				isAdmin = true
			}
		}
	}

	cmd.IsAdmin = isAdmin

	return cmd, nil
}

// IsAdminCommand returns true when it is a admin command, otherwise false.
func (cmd *Command) IsAdminCommand() bool {
	return cmd.IsAdmin
}

// GetType returns a string type
func (cmd *Command) GetType() string {
	return cmd.Type
}

// IsType returns true when the command has the specified element, otherwise false.
func (cmd *Command) IsType(typeString string) bool {
	if cmd.Type != typeString {
		return false
	}
	return true
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
