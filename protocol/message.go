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
	"fmt"
)

// Message represents a message of MongoDB Wire
// See : MongoDB Wire Protocol
// https://docs.mongodb.com/manual/reference/mongodb-wire-protocol/

// MessageListener represents a listener for MongoDB Wire Protocol.
type MessageListener interface {
	MessageReceived(Message)
	MessageRespond(Message)
}

// Message represents an operation message of MongoDB wire protocol.
type Message interface {
	// SetRequestID sets a message identifier.
	SetRequestID(id int32)
	// SetResponseTo sets a response message identifier.
	SetResponseTo(id int32)
	// GetMessageLength returns the message length.
	GetMessageLength() int32
	// GetRequestID returns the message identifier.
	GetRequestID() int32
	// GetResponseTo returns the response message identifier.
	GetResponseTo() int32
	// GetOpCode returns the operation code.
	GetOpCode() OpCode
	// Size returns the message size including the header.
	Size() int32
	// Bytes returns the binary description of BSON format.
	Bytes() []byte
	// String returns the string description.
	String() string
}

// NewMessageWithHeaderAndBytes returns a parsed message of the specified header and body bytes.
func NewMessageWithHeaderAndBytes(header *Header, body []byte) (Message, error) {
	switch header.OpCode {
	case OpUpdate:
		return NewUpdateWithHeaderAndBody(header, body)
	case OpInsert:
		return NewInsertWithHeaderAndBody(header, body)
	case OpQuery:
		return NewQueryWithHeaderAndBody(header, body)
	case OpGetMore:
		return NewGetMoreWithHeaderAndBody(header, body)
	case OpDelete:
		return NewDeleteWithHeaderAndBody(header, body)
	case OpKillCursors:
		return NewKillCursorsWithHeaderAndBody(header, body)
	case OpMsg:
		return NewMsgWithHeaderAndBody(header, body)
	case OpReply:
		return NewReplyWithHeaderAndBody(header, body)
	default:
	}
	return nil, fmt.Errorf(errorInvalidMessageOpCode, header.OpCode)
}

// NewMessageWithBytes returns a parsed message of the specified bytes.
func NewMessageWithBytes(msg []byte) (Message, error) {
	header, err := NewHeaderWithBytes(msg)
	if err != nil {
		return nil, err
	}
	return NewMessageWithHeaderAndBytes(header, msg[:HeaderSize])
}
