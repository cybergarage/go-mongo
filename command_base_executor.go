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

	"github.com/cybergarage/go-mongo/bson"
	"github.com/cybergarage/go-mongo/message"
)

// BaseCommandExecutor is a complete hander for CommandExecutor.
type BaseCommandExecutor struct {
	DatabaseCommandExecutor
}

func baseCommandExecutorNotImplementedError(q *Query) error {
	return fmt.Errorf(errorQueryHanderNotImplemented, "")
}

// NewBaseCommandExecutor returns a complete null executor for CommandExecutor.
func NewBaseCommandExecutor() *BaseCommandExecutor {
	executor := &BaseCommandExecutor{}
	executor.DatabaseCommandExecutor = executor
	return executor
}

// SetDatabaseCommandExecutor sets a query exector for database operation commands.
func (executor *BaseCommandExecutor) SetDatabaseCommandExecutor(fn DatabaseCommandExecutor) {
	executor.DatabaseCommandExecutor = fn
}

//////////////////////////////////////////////////
// CommandExecutor
//////////////////////////////////////////////////

// ExecuteCommand handles query commands other than those explicitly specified above.
func (executor *BaseCommandExecutor) ExecuteCommand(cmd *Command) ([]bson.Document, error) {
	if executor.DatabaseCommandExecutor == nil {
		return nil, fmt.Errorf(errorQueryHanderNotImplemented, cmd.String())
	}

	cmdType, err := cmd.GetType()
	if err != nil {
		return nil, err
	}

	switch cmdType {
	case message.IsMaster:
		return executor.DatabaseCommandExecutor.ExecuteIsMaster(cmd)
	}

	return nil, fmt.Errorf(errorQueryHanderNotImplemented, cmd.String())
}

// IsMaster displays information about this member’s role in the replica set, including whether it is the master.
func (executor *BaseCommandExecutor) ExecuteIsMaster(cmd *Command) ([]bson.Document, error) {
	reply := message.NewDefaultIsMasterResponse()
	replyDoc, err := reply.BSONBytes()
	if err != nil {
		return nil, err
	}
	return []bson.Document{replyDoc}, nil
}

//////////////////////////////////////////////////
// CommandExecutor
//////////////////////////////////////////////////

// Insert hadles OP_INSERT and 'insert' query of OP_MSG.
func (executor *BaseCommandExecutor) Insert(*Query) (int32, bool) {
	return 0, false
}

// Update hadles OP_UPDATE and 'update' query of OP_MSG.
func (executor *BaseCommandExecutor) Update(*Query) (int32, bool) {
	return 0, false
}

// Find hadles 'find' query of OP_MSG.
func (executor *BaseCommandExecutor) Find(*Query) ([]bson.Document, bool) {
	return nil, false
}

// Delete hadles OP_DELETE and 'delete' query of OP_MSG.
func (executor *BaseCommandExecutor) Delete(*Query) (int32, bool) {
	return 0, false
}
