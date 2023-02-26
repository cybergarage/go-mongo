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
	"github.com/cybergarage/go-mongo/mongo/protocol"
)

// OpReply (OP_REPLY) replies to a client request. responseTo is set.
type OpReply = protocol.Reply

// OpUpdate (OP_UPDATE:2001) updates document.
type OpUpdate = protocol.Update

// OpInsert (OP_INSERT:2002) inserts new document.
type OpInsert = protocol.Insert

// OpQuery (OP_QUERY:2004) queries a collection.
type OpQuery = protocol.Query

// OpGetMore (GET_MORE:2005) gets more data from a query. See Cursors.
type OpGetMore = protocol.GetMore

// OpDelete (OP_DELETE:2006) deletes documents.
type OpDelete = protocol.Delete

// OpKillCursors (OP_KILL_CURSORS:2007) notifies database that the client has finished with the cursor.
type OpKillCursors = protocol.KillCursors

// OpMsg (OP_MSG:2013) sends a message using the format introduced in MongoDB 3.6.
type OpMsg = protocol.Msg

// OpFlag represents a flag of MongoDB wire protocol.
type OpFlag = protocol.Flag

// OpMessage represents a message of MongoDB wire protocol.
type OpMessage = protocol.Message

////////////////////////////////////////
// OpMessageHandler
////////////////////////////////////////

// OpMessageHandler represents an interface for MongoDB query request.
type OpMessageHandler interface {
	// Update handles OP_UPDATE of MongoDB wire protocol.
	OpUpdate(conn *Conn, q *OpUpdate) (bson.Document, error)
	// Insert handles OP_INSERT of MongoDB wire protocol.
	OpInsert(conn *Conn, q *OpInsert) (bson.Document, error)
	// Query handles OP_QUERY of MongoDB wire protocol.
	OpQuery(conn *Conn, q *OpQuery) (bson.Document, error)
	// GetMore handles GET_MORE of MongoDB wire protocol.
	OpGetMore(conn *Conn, q *OpGetMore) (bson.Document, error)
	// Delete handles OP_DELETE of MongoDB wire protocol.
	OpDelete(conn *Conn, q *OpDelete) (bson.Document, error)
	// KillCursors handles OP_KILL_CURSORS of MongoDB wire protocol.
	OpKillCursors(conn *Conn, q *OpKillCursors) (bson.Document, error)
	// Msg handles OP_MSG of MongoDB wire protocol.
	OpMsg(conn *Conn, q *OpMsg) (bson.Document, error)
}
