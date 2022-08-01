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

package protocol

import (
	"go.mongodb.org/mongo-driver/x/mongo/driver/wiremessage"
)

// OpCode represents a MongoDB wire protocol opcode.
type OpCode = wiremessage.OpCode

const (
	//OpReply (OP_REPLY) replies to a client request. responseTo is set.
	OpReply OpCode = wiremessage.OpReply
	// OpUpdate (OP_UPDATE:2001) updates document.
	OpUpdate OpCode = wiremessage.OpUpdate
	// OpInsert (OP_INSERT:2002) inserts new document.
	OpInsert OpCode = wiremessage.OpInsert
	// OpQuery (OP_QUERY:2004) queries a collection.
	OpQuery OpCode = wiremessage.OpQuery
	// OpGetMore (GET_MORE:2005) gets more data from a query. See Cursors.
	OpGetMore OpCode = wiremessage.OpGetMore
	// OpDelete (OP_DELETE:2006) Deletes documents.
	OpDelete OpCode = wiremessage.OpDelete
	// OpKillCursors (OP_KILL_CURSORS:2007)Notifies database that the client has finished with the cursor.
	OpKillCursors OpCode = wiremessage.OpKillCursors
	// OpCommand (OP_COMMAND:2011) clusters internal protocol representing a command request.
	OpCommand OpCode = wiremessage.OpCommand
	// OpCommandReply (OP_COMMANDREPLY:2011) clusters internal protocol representing a reply to an OP_COMMAND.
	OpCommandReply OpCode = wiremessage.OpCommandReply
	// OpMsg (OP_MSG:2013) sends a message using the format introduced in MongoDB 3.6.
	OpMsg OpCode = wiremessage.OpMsg
)
