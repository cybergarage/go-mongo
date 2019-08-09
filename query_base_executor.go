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

// BaseQueryExecutor is a complete hander for QueryExecutor.
type BaseQueryExecutor struct{}

func baseQueryExecutorNotImplementedError(q *Query) error {
	return fmt.Errorf(errorQueryHanderNotImplemented, "")
}

// NewBaseQueryExecutor returns a complete null handler for QueryExecutor.
func NewBaseQueryExecutor() *BaseQueryExecutor {
	return &BaseQueryExecutor{}
}

//////////////////////////////////////////////////
// CommandExecutor
//////////////////////////////////////////////////

// ExecuteCommand handles query commands other than those explicitly specified above.
func (handler *BaseQueryExecutor) ExecuteCommand(cmd *Command) ([]bson.Document, error) {
	cmdType, err := cmd.GetType()
	if err != nil {
		return nil, err
	}
	switch cmdType {
	case message.IsMaster:
		return handler.isMaster(cmd)
	}
	return nil, fmt.Errorf(errorQueryHanderNotImplemented, cmd.String())
}

// Replication Commands

type ReplicationExecutor interface {
	// isMaster displays information about this member’s role in the replica set, including whether it is the master.
	isMaster(q *Query) ([]bson.Document, error)
}

// IsMaster displays information about this member’s role in the replica set, including whether it is the master.
func (handler *BaseQueryExecutor) isMaster(cmd *Command) ([]bson.Document, error) {
	reply := message.NewDefaultIsMasterResponse()
	replyDoc, err := reply.BSONBytes()
	if err != nil {
		return nil, err
	}
	return []bson.Document{replyDoc}, nil
}

//////////////////////////////////////////////////
// QueryExecutor
//////////////////////////////////////////////////

// Insert hadles OP_INSERT and 'insert' query of OP_MSG.
func (handler *BaseQueryExecutor) Insert(*Query) (int32, bool) {
	return 0, false
}

// Update hadles OP_UPDATE and 'update' query of OP_MSG.
func (handler *BaseQueryExecutor) Update(*Query) (int32, bool) {
	return 0, false
}

// Find hadles 'find' query of OP_MSG.
func (handler *BaseQueryExecutor) Find(*Query) ([]bson.Document, bool) {
	return nil, false
}

// Delete hadles OP_DELETE and 'delete' query of OP_MSG.
func (handler *BaseQueryExecutor) Delete(*Query) (int32, bool) {
	return 0, false
}
