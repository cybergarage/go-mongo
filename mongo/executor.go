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
	"github.com/cybergarage/go-mongo/mongo/bson"
	"github.com/cybergarage/go-mongo/mongo/message"
)

// Query represents a query of MongoDB database command.
type Query = message.Query

// Executor represents an executor interface for MongoDB message commands.
type Executor interface {
	QueryCommandExecutor
}

// QueryCommandExecutor represents an executor interface for MongoDB queries.
type QueryCommandExecutor interface {
	// Insert hadles OP_INSERT and 'insert' query of OP_MSG.
	Insert(*Conn, *Query) (int32, error)
	// Update hadles OP_UPDATE and 'update' query of OP_MSG.
	Update(*Conn, *Query) (int32, error)
	// Find hadles 'find' query of OP_MSG.
	Find(*Conn, *Query) ([]bson.Document, error)
	// Delete hadles OP_DELETE and 'delete' query of OP_MSG.
	Delete(*Conn, *Query) (int32, error)
}
