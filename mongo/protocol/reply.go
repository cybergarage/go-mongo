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

	"github.com/cybergarage/go-mongo/mongo/bson"
	"go.mongodb.org/mongo-driver/x/mongo/driver/wiremessage"
)

const (

	// CursorNotFound sets when getMore is called but the cursor id is not valid at the server. Returned with zero results.
	CursorNotFound = 0x01
	// QueryFailure sets when query failed. Results consist of one document containing an “$err” field describing the failure.
	QueryFailure = 0x02
	// ShardConfigStale needs to update config from the server, and so drivers should ignore this.
	ShardConfigStale = 0x04
	// AwaitCapable sets when the server supports the AwaitData Query option. If it doesn’t, a client should sleep a little between getMore’s of a Tailable cursor. Mongod version 1.6 supports AwaitData and thus always sets AwaitCapable.
	AwaitCapable = 0x08
)

// ReplyFlag represents a MsgFReplyFlaglag of MongoDB wire protocol.
type ReplyFlag = wiremessage.ReplyFlag

// Reply represents a OP_REPLY of MongoDB wire protocol.
// See : MongoDB Wire Protocol
// https://docs.mongodb.com/manual/reference/mongodb-wire-protocol/
type Reply struct {
	*Header // A standard wire protocol header

	ReplyFlags     ReplyFlag       // bit vector - see details below
	CursorID       int64           // cursor id if client needs to do get more's
	StartingFrom   int32           // where in the cursor this reply is starting
	NumberReturned int32           // number of documents in the reply
	documents      []bson.Document // documents
}

// NewReply returns a new reply instance.
func NewReply() *Reply {
	op := &Reply{
		Header:         NewHeaderWithOpCode(OpReply),
		ReplyFlags:     0,
		CursorID:       0,
		StartingFrom:   0,
		NumberReturned: 0,
		documents:      make([]bson.Document, 0),
	}

	return op
}

// NewReplyWithDocuments returns a new reply instance with ths specified documents.
func NewReplyWithDocuments(docs []bson.Document) *Reply {
	op := NewReply()
	op.NumberReturned = int32(len(docs))
	op.documents = docs
	op.SetMessageLength(op.Size())
	return op
}

// NewReplyWithDocument returns a new reply instance with ths specified document.
func NewReplyWithDocument(doc bson.Document) *Reply {
	return NewReplyWithDocuments([]bson.Document{doc})
}

// NewReplyWithHeaderAndBody returns a new reply instance with the specified bytes.
func NewReplyWithHeaderAndBody(header *Header, body []byte) (*Reply, error) {
	responseFlags, offsetBody, ok := ReadInt32(body)
	if !ok {
		return nil, newErrMessageRequest(OpReply, body)
	}

	cursorID, offsetBody, ok := ReadInt64(offsetBody)
	if !ok {
		return nil, newErrMessageRequest(OpReply, body)
	}

	startingFrom, offsetBody, ok := ReadInt32(offsetBody)
	if !ok {
		return nil, newErrMessageRequest(OpReply, body)
	}

	numberReturned, offsetBody, ok := ReadInt32(offsetBody)
	if !ok {
		return nil, newErrMessageRequest(OpReply, body)
	}

	documents, _, ok := ReadDocuments(offsetBody)
	if !ok {
		return nil, newErrMessageRequest(OpReply, body)
	}

	op := &Reply{
		Header:         header,
		ReplyFlags:     ReplyFlag(responseFlags),
		CursorID:       cursorID,
		StartingFrom:   startingFrom,
		NumberReturned: numberReturned,
		documents:      documents,
	}

	return op, nil
}

// SetResponseFlags sets a response flag.
func (op *Reply) SetResponseFlags(flag ReplyFlag) {
	op.ReplyFlags = flag
}

// Documents returns the BSON documents.
func (op *Reply) Documents() []bson.Document {
	return op.documents
}

// Size returns the message size including the header.
func (op *Reply) Size() int32 {
	bodySize := 4 + 8 + 4 + 4
	for _, document := range op.documents {
		bodySize += len(document)
	}
	return int32(HeaderSize + bodySize)
}

// Bytes returns the binary description of BSON format.
func (op *Reply) Bytes() []byte {
	dst := op.Header.Bytes()
	dst = AppendInt32(dst, int32(op.ReplyFlags))
	dst = AppendInt64(dst, op.CursorID)
	dst = AppendInt32(dst, op.StartingFrom)
	dst = AppendInt32(dst, op.NumberReturned)
	for _, document := range op.documents {
		dst = AppendDocument(dst, document)
	}
	return dst
}

// String returns the string description.
func (op *Reply) String() string {
	str := fmt.Sprintf("%s %d %d %d %d ",
		op.Header.String(),
		op.ReplyFlags,
		op.CursorID,
		op.StartingFrom,
		op.NumberReturned,
	)

	for _, document := range op.documents {
		str += fmt.Sprintf("%s ", document.String())
	}

	return str
}
