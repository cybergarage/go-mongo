// Copyright (C) 2019 The go-mongo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mongo

import (
	"github.com/cybergarage/go-mongo/mongo/bson"
	"github.com/cybergarage/go-mongo/mongo/message"
)

// Command represents a query command of MongoDB database command.
type Command = message.Command

// Query represents a query of MongoDB database command.
type Query = message.Query

// CommandExecutor represents an executor interface for MongoDB commands.
type CommandExecutor = message.CommandExecutor

// MessageExecutor represents an executor interface for MongoDB message.
type MessageExecutor = message.MessageExecutor

// UserCommandExecutor represents an executor interface for MongoDB query commands.
type UserCommandExecutor interface {
	message.QueryCommandExecutor
}

// DatabaseCommandExecutor represents an executor interface for MongoDB operation commands.
type DatabaseCommandExecutor interface {
	ReplicationCommandExecutor
	DiagnosticCommandExecutor
	WriteOperationExecutor
}

// ReplicationCommandExecutor represents an executor interface for MongoDB replication commands.
type ReplicationCommandExecutor interface {
	ExecuteIsMaster(cmd *Command) (bson.Document, error)
}

// DiagnosticCommandExecutor represents an executor interface for MongoDB diagnostic commands.
type DiagnosticCommandExecutor interface {
	ExecuteBuildInfo(cmd *Command) (bson.Document, error)
}

// WriteOperationExecutor represents an executor interface for MongoDB write operation commands.
type WriteOperationExecutor interface {
	ExecuteGetLastError(cmd *Command) (bson.Document, error)
}
