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

// BaseMessageHandler is a complete hander for MessageHandler.
type BaseMessageHandler struct {
	CommandExecutor
	MessageExecutor
}

func newBaseMessageHandlerNotImplementedError(msg OpMessage) error {
	return fmt.Errorf(errorMessageHanderNotSupported, msg.OpCode())
}

// NewBaseMessageHandler returns a complete null handler for MessageHandler.
func NewBaseMessageHandler() *BaseMessageHandler {
	return &BaseMessageHandler{
		CommandExecutor: nil,
		MessageExecutor: nil,
	}
}

// SetCommandExecutor sets a exector for OP_QUERY of MongoDB wire protocol.
func (handler *BaseMessageHandler) SetCommandExecutor(fn CommandExecutor) {
	handler.CommandExecutor = fn
}

// SetMessageExecutor sets a exector for OP_MSG of MongoDB wire protocol.
func (handler *BaseMessageHandler) SetMessageExecutor(fn MessageExecutor) {
	handler.MessageExecutor = fn
}

// OpUpdate handles OP_UPDATE of MongoDB wire protocol.
func (handler *BaseMessageHandler) OpUpdate(conn *Conn, msg *OpUpdate) (bson.Document, error) {
	return nil, newBaseMessageHandlerNotImplementedError(msg)
}

// OpInsert handles OP_INSERT of MongoDB wire protocol.
func (handler *BaseMessageHandler) OpInsert(conn *Conn, msg *OpInsert) (bson.Document, error) {
	return nil, newBaseMessageHandlerNotImplementedError(msg)
}

// OpQuery handles OP_QUERY of MongoDB wire protocol.
func (handler *BaseMessageHandler) OpQuery(conn *Conn, msg *OpQuery) (bson.Document, error) {
	if handler.CommandExecutor == nil {
		return nil, newBaseMessageHandlerNotImplementedError(msg)
	}

	cmd, err := message.NewCommandWithQuery(msg)
	if err != nil {
		return nil, err
	}

	cmdType := cmd.GetType()
	conn.StartSpan(cmdType)
	defer conn.FinishSpan()

	switch cmdType {
	// For user database commands over OP_QUERY under MongoDB v3.6
	case message.Insert, message.Delete, message.Update, message.Find:
		q, err := message.NewQueryWithQuery(msg)
		if err != nil {
			return nil, err
		}
		res := message.NewResponse()
		err = handler.executeQuery(conn, q, res)
		if err != nil {
			return nil, err
		}
		bsonRes, err := res.BSONBytes()
		if err != nil {
			return nil, err
		}
		return bsonRes, nil
	}

	return handler.CommandExecutor.ExecuteCommand(conn, cmd)
}

// OpGetMore handles GET_MORE of MongoDB wire protocol.
func (handler *BaseMessageHandler) OpGetMore(conn *Conn, msg *OpGetMore) (bson.Document, error) {
	return nil, newBaseMessageHandlerNotImplementedError(msg)
}

// OpDelete handles OP_DELETE of MongoDB wire protocol.
func (handler *BaseMessageHandler) OpDelete(conn *Conn, msg *OpDelete) (bson.Document, error) {
	return nil, newBaseMessageHandlerNotImplementedError(msg)
}

// OpKillCursors handles OP_KILL_CURSORS of MongoDB wire protocol.
func (handler *BaseMessageHandler) OpKillCursors(conn *Conn, msg *OpKillCursors) (bson.Document, error) {
	// TODO : Kill the specified cursors internally
	res := message.NewResponse()
	res.SetStatus(true)
	bsonRes, err := res.BSONBytes()
	if err != nil {
		return nil, err
	}
	return bsonRes, nil
}

// OpMsg handles OP_MSG of MongoDB wire protocol.
func (handler *BaseMessageHandler) OpMsg(conn *Conn, msg *OpMsg) (bson.Document, error) {
	if handler.MessageExecutor == nil {
		return nil, newBaseMessageHandlerNotImplementedError(msg)
	}

	q, err := message.NewQueryWithMessage(msg)
	if err != nil {
		return nil, err
	}

	res := message.NewResponse()

	queryType := q.GetType()
	switch queryType {
	// For user database commands over OP_MSG from MongoDB v3.6
	case message.Insert, message.Delete, message.Update, message.Find:
		conn.StartSpan(queryType)
		defer conn.FinishSpan()
		err = handler.executeQuery(conn, q, res)
		if err != nil {
			return nil, err
		}
	case message.KillCursors:
		conn.StartSpan(queryType)
		defer conn.FinishSpan()
		// TODO : Kill the specified cursors internally
		res.SetStatus(true)
	default: // Execute other messages as a database command
		cmd, err := message.NewCommandWithMsg(msg)
		if err != nil {
			return nil, err
		}
		conn.StartSpan(cmd.String())
		defer conn.FinishSpan()
		resDoc, err := handler.CommandExecutor.ExecuteCommand(conn, cmd)
		if err != nil {
			return nil, err
		}
		return resDoc, nil
	}

	bsonRes, err := res.BSONBytes()
	if err != nil {
		return nil, err
	}

	return bsonRes, nil
}

// executeQuery executes user database commands (insert, update, find and delete) over OP_MSG and OP_QUERY.
func (handler *BaseMessageHandler) executeQuery(conn *Conn, q *message.Query, res *message.Response) error {
	switch q.GetType() {
	case message.Insert:
		n, err := handler.MessageExecutor.Insert(conn, q)
		res.SetErrorStatus(err)
		res.SetNumberOfAffectedDocuments(n)
	case message.Delete:
		n, err := handler.MessageExecutor.Delete(conn, q)
		res.SetErrorStatus(err)
		res.SetNumberOfAffectedDocuments(n)
	case message.Update:
		n, err := handler.MessageExecutor.Update(conn, q)
		res.SetErrorStatus(err)
		res.SetNumberOfAffectedDocuments(n)
		res.SetNumberOfModifiedDocuments(n)
	case message.Find:
		docs, err := handler.MessageExecutor.Find(conn, q)
		res.SetErrorStatus(err)
		res.SetCursorDocuments(q.GetFullCollectionName(), docs)
	default:
		res.SetStatus(false)
	}
	return nil
}
