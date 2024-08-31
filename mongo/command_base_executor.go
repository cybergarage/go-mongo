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

package mongo

import (
	"fmt"

	"github.com/cybergarage/go-mongo/mongo/bson"
	"github.com/cybergarage/go-mongo/mongo/message"
)

// BaseCommandExecutor is a complete hander for CommandExecutor.
type BaseCommandExecutor struct {
	UserCommandExecutor
	DatabaseCommandExecutor
	AuthCommandExecutor
}

func baseCommandExecutorNotImplementedError(q *Query) error {
	return fmt.Errorf(errorQueryHanderNotImplemented, "")
}

// NewBaseCommandExecutor returns a complete null executor for CommandExecutor.
func NewBaseCommandExecutor() *BaseCommandExecutor {
	executor := &BaseCommandExecutor{
		UserCommandExecutor:     nil,
		DatabaseCommandExecutor: nil,
		AuthCommandExecutor:     nil,
	}
	executor.UserCommandExecutor = executor
	executor.DatabaseCommandExecutor = executor
	executor.AuthCommandExecutor = executor
	return executor
}

// SetUserCommandExecutor sets a command exector for database operation commands.
func (executor *BaseCommandExecutor) SetUserCommandExecutor(fn UserCommandExecutor) {
	executor.UserCommandExecutor = fn
}

// SetDatabaseCommandExecutor sets a command exector for database operation commands.
func (executor *BaseCommandExecutor) SetDatabaseCommandExecutor(fn DatabaseCommandExecutor) {
	executor.DatabaseCommandExecutor = fn
}

// SetAuthCommandExecutor  sets a command exector for auth operation commands.
func (executor *BaseCommandExecutor) SetAuthCommandExecutor(fn AuthCommandExecutor) {
	executor.AuthCommandExecutor = fn
}

//////////////////////////////////////////////////
// CommandExecutor
//////////////////////////////////////////////////

// ExecuteCommand handles query commands other than those explicitly specified above.
func (executor *BaseCommandExecutor) ExecuteCommand(conn *Conn, cmd *Command) (bson.Document, error) {
	if executor.DatabaseCommandExecutor == nil {
		// Returns only a 'ok' response as default
		resDoc, err := message.NewOkResponse().BSONBytes()
		if err != nil {
			return nil, err
		}
		return resDoc, nil
	}

	switch cmd.Type() {
	case message.IsMaster:
		return executor.DatabaseCommandExecutor.Hello(conn, cmd)
	case message.BuildInfo:
		return executor.DatabaseCommandExecutor.BuildInfo(conn, cmd)
	case message.GetLastError:
		return executor.DatabaseCommandExecutor.GetLastError(conn, cmd)
	case message.SASLStart:
		return executor.AuthCommandExecutor.SASLStart(conn, cmd)
	case message.SASLContinue:
		return executor.AuthCommandExecutor.SASLContinue(conn, cmd)
	}

	// Returns only a 'ok' response as default
	resDoc, err := message.NewOkResponse().BSONBytes()
	if err != nil {
		return nil, err
	}
	return resDoc, nil
}

//////////////////////////////////////////////////
// DatabaseCommandExecutor
//////////////////////////////////////////////////

// Hello returns information about this member’s role in the replica set, including whether it is the master.
func (executor *BaseCommandExecutor) Hello(conn *Conn, cmd *Command) (bson.Document, error) {
	reply, err := message.NewDefaultIsMasterResponse()
	if err != nil {
		return nil, err
	}
	replyDoc, err := reply.BSONBytes()
	if err != nil {
		return nil, err
	}
	return replyDoc, nil
}

// BuildInfo returns statistics about the MongoDB build.
func (executor *BaseCommandExecutor) BuildInfo(conn *Conn, cmd *Command) (bson.Document, error) {
	reply, err := message.NewDefaultBuildInfoResponse()
	if err != nil {
		return nil, err
	}
	replyDoc, err := reply.BSONBytes()
	if err != nil {
		return nil, err
	}
	return replyDoc, nil
}

// GetLastError returns statistics about the MongoDB build.
func (executor *BaseCommandExecutor) GetLastError(conn *Conn, cmd *Command) (bson.Document, error) {
	reply, err := message.NewDefaultLastErrorResponse()
	if err != nil {
		return nil, err
	}
	replyDoc, err := reply.BSONBytes()
	if err != nil {
		return nil, err
	}
	return replyDoc, nil
}

//////////////////////////////////////////////////
// UserCommandExecutor
//////////////////////////////////////////////////

// Insert hadles OP_INSERT and 'insert' query of OP_MSG.
func (executor *BaseCommandExecutor) Insert(conn *Conn, q *Query) (int32, error) {
	if executor.UserCommandExecutor != nil {
		return executor.UserCommandExecutor.Insert(conn, q)
	}
	return 0, NewNotSupported(q)
}

// Update hadles OP_UPDATE and 'update' query of OP_MSG.
func (executor *BaseCommandExecutor) Update(conn *Conn, q *Query) (int32, error) {
	if executor.UserCommandExecutor != nil {
		return executor.UserCommandExecutor.Update(conn, q)
	}
	return 0, NewNotSupported(q)
}

// Find hadles 'find' query of OP_MSG.
func (executor *BaseCommandExecutor) Find(conn *Conn, q *Query) ([]bson.Document, error) {
	if executor.UserCommandExecutor != nil {
		return executor.UserCommandExecutor.Find(conn, q)
	}
	return nil, NewNotSupported(q)
}

// Delete hadles OP_DELETE and 'delete' query of OP_MSG.
func (executor *BaseCommandExecutor) Delete(conn *Conn, q *Query) (int32, error) {
	if executor.UserCommandExecutor != nil {
		return executor.UserCommandExecutor.Delete(conn, q)
	}
	return 0, NewNotSupported(q)
}

//////////////////////////////////////////////////
// AuthCommandExecutor
//////////////////////////////////////////////////

// SASLStart handles SASLStart command.
func (executor *BaseCommandExecutor) SASLStart(conn *Conn, cmd *Command) (bson.Document, error) {
	return nil, nil
}

// SASLContinue handles SASLContinue command.
func (executor *BaseCommandExecutor) SASLContinue(*Conn, *Command) (bson.Document, error) {
	return nil, nil
}
