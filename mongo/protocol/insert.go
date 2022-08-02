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

// Insert represents a OP_INSERT of MongoDB wire protocol.
// See : MongoDB Wire Protocol
// https://docs.mongodb.com/manual/reference/mongodb-wire-protocol/
type Insert struct {
	*Header                          // A standard wire protocol header
	Flags              Flag          // bit vector. see below
	FullCollectionName string        // "dbname.collectionname"
	Document           bson.Document // one or more documents to insert into the collection
}

// NewInsertWithHeaderAndBody returns a new insert instance with the specified bytes.
func NewInsertWithHeaderAndBody(header *Header, body []byte) (*Insert, error) {
	flags, offsetBody, ok := ReadUint32(body)
	if !ok {
		return nil, newMessageRequestError(OpInsert, body)
	}

	collectionName, offsetBody, ok := ReadCString(offsetBody)
	if !ok {
		return nil, newMessageRequestError(OpInsert, body)
	}

	document, _, ok := ReadDocument(offsetBody)
	if !ok {
		return nil, newMessageRequestError(OpInsert, body)
	}

	op := &Insert{
		Header:             header,
		Flags:              Flag(flags),
		FullCollectionName: collectionName,
		Document:           document,
	}

	return op, nil
}

// Size returns the message size including the header.
func (op *Insert) Size() int32 {
	bodySize := 4 + (len(op.FullCollectionName) + 1) + len(op.Document)
	return int32(HeaderSize + bodySize)
}

// String returns the string description.
func (op *Insert) String() string {
	return fmt.Sprintf("%s %X %s %s",
		op.Header.String(),
		op.Flags,
		op.FullCollectionName,
		op.Document.String(),
	)
}
