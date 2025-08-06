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
)

// Update represents a OP_UPDATE of MongoDB wire protocol.
// See : MongoDB Wire Protocol
// https://docs.mongodb.com/manual/reference/mongodb-wire-protocol/
type Update struct {
	*Header // A standard wire protocol header

	ZERO               int32         // 0 - reserved for future use
	FullCollectionName string        // "dbname.collectionname"
	Flags              Flag          // bit vector. see below
	Selector           bson.Document // the query to select the document
	Update             bson.Document // specification of the update to perform
}

// NewUpdateWithHeaderAndBody returns a new update instance with the specified bytes.
func NewUpdateWithHeaderAndBody(header *Header, body []byte) (*Update, error) {
	zero, offsetBody, ok := ReadInt32(body)
	if !ok {
		return nil, newErrMessageRequest(OpDelete, body)
	}

	collectionName, offsetBody, ok := ReadCString(offsetBody)
	if !ok {
		return nil, newErrMessageRequest(OpUpdate, body)
	}

	flags, offsetBody, ok := ReadUint32(offsetBody)
	if !ok {
		return nil, newErrMessageRequest(OpUpdate, body)
	}

	selector, offsetBody, ok := ReadDocument(offsetBody)
	if !ok {
		return nil, newErrMessageRequest(OpUpdate, body)
	}

	update, _, ok := ReadDocument(offsetBody)
	if !ok {
		return nil, newErrMessageRequest(OpUpdate, body)
	}

	op := &Update{
		Header:             header,
		ZERO:               zero,
		FullCollectionName: collectionName,
		Selector:           selector,
		Update:             update,
		Flags:              Flag(flags),
	}

	return op, nil
}

// Documents returns the BSON documents.
func (op *Update) Documents() []bson.Document {
	return []bson.Document{op.Selector, op.Update}
}

// Size returns the message size including the header.
func (op *Update) Size() int32 {
	bodySize := 4 + (len(op.FullCollectionName) + 1) + 4 + len(op.Selector) + len(op.Update)
	return int32(HeaderSize + bodySize)
}

// String returns the string description.
func (op *Update) String() string {
	return fmt.Sprintf("%s %s %X %s %s",
		op.Header.String(),
		op.FullCollectionName,
		op.Flags,
		op.Selector.String(),
		op.Update.String(),
	)
}
