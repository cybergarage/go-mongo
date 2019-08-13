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
	return fmt.Errorf(errorMessageHanderNotSupported, msg.GetOpCode())
}

// NewBaseMessageHandler returns a complete null handler for MessageHandler.
func NewBaseMessageHandler() *BaseMessageHandler {
	return &BaseMessageHandler{
		CommandExecutor: nil,
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
func (handler *BaseMessageHandler) OpUpdate(msg *OpUpdate) (bson.Document, error) {
	return nil, newBaseMessageHandlerNotImplementedError(msg)
}

// OpInsert handles OP_INSERT of MongoDB wire protocol.
func (handler *BaseMessageHandler) OpInsert(msg *OpInsert) (bson.Document, error) {
	return nil, newBaseMessageHandlerNotImplementedError(msg)
}

// OpQuery handles OP_QUERY of MongoDB wire protocol.
func (handler *BaseMessageHandler) OpQuery(msg *OpQuery) ([]bson.Document, error) {
	if handler.CommandExecutor == nil {
		return nil, newBaseMessageHandlerNotImplementedError(msg)
	}

	cmd, err := message.NewCommandWithQuery(msg)
	if err != nil {
		return nil, err
	}

	return handler.CommandExecutor.ExecuteCommand(cmd)
}

// OpGetMore handles GET_MORE of MongoDB wire protocol.
func (handler *BaseMessageHandler) OpGetMore(msg *OpGetMore) (bson.Document, error) {
	return nil, newBaseMessageHandlerNotImplementedError(msg)
}

// OpDelete handles OP_DELETE of MongoDB wire protocol.
func (handler *BaseMessageHandler) OpDelete(msg *OpDelete) (bson.Document, error) {
	return nil, newBaseMessageHandlerNotImplementedError(msg)
}

// OpKillCursors handles OP_KILL_CURSORS of MongoDB wire protocol.
func (handler *BaseMessageHandler) OpKillCursors(msg *OpKillCursors) (bson.Document, error) {
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
func (handler *BaseMessageHandler) OpMsg(msg *OpMsg) (bson.Document, error) {
	if handler.MessageExecutor == nil {
		return nil, newBaseMessageHandlerNotImplementedError(msg)
	}

	q, err := message.NewQueryWithMessage(msg)
	if err != nil {
		return nil, err
	}

	res := message.NewResponse()

	switch q.GetType() {
	case message.Insert:
		n, ok := handler.MessageExecutor.Insert(q)
		res.SetStatus(ok)
		res.SetNumberOfAffectedDocuments(n)
	case message.Delete:
		n, ok := handler.MessageExecutor.Delete(q)
		res.SetStatus(ok)
		res.SetNumberOfAffectedDocuments(n)
	case message.Update:
		n, ok := handler.MessageExecutor.Update(q)
		res.SetStatus(ok)
		res.SetNumberOfAffectedDocuments(n)
		res.SetNumberOfModifiedDocuments(n)
	case message.Find:
		docs, ok := handler.MessageExecutor.Find(q)
		res.SetCursorDocuments(q.GetFullCollectionName(), docs)
		res.SetStatus(ok)
	case message.KillCursors:
		// TODO : Kill the specified cursors internally
		res.SetStatus(true)
	}

	bsonRes, err := res.BSONBytes()
	if err != nil {
		return nil, err
	}

	return bsonRes, nil
}
