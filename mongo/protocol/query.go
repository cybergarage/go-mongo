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

// Query represents a OP_QUERY of MongoDB wire protocol.
// See : MongoDB Wire Protocol
// https://docs.mongodb.com/manual/reference/mongodb-wire-protocol/
type Query struct {
	*Header                          // A standard wire protocol header
	Flags              Flag          // bit vector. see below
	FullCollectionName string        // "dbname.collectionname"
	NumberToSkip       int32         // number of documents to skip
	NumberToReturn     int32         // number of documents to return in the first OP_REPLY batch
	Query              bson.Document // query object.  See below for details.
}

// NewQueryWithHeaderAndBody returns a new insert instance with the specified bytes.
func NewQueryWithHeaderAndBody(header *Header, body []byte) (*Query, error) {
	flags, offsetBody, ok := ReadUint32(body)
	if !ok {
		return nil, newMessageRequestError(OpQuery, body)
	}

	collectionName, offsetBody, ok := ReadCString(offsetBody)
	if !ok {
		return nil, newMessageRequestError(OpQuery, body)
	}

	numberToSkip, offsetBody, ok := ReadInt32(offsetBody)
	if !ok {
		return nil, newMessageRequestError(OpQuery, body)
	}

	numberToReturn, offsetBody, ok := ReadInt32(offsetBody)
	if !ok {
		return nil, newMessageRequestError(OpQuery, body)
	}

	query, _, ok := ReadDocument(offsetBody)
	if !ok {
		return nil, newMessageRequestError(OpInsert, body)
	}

	op := &Query{
		Header:             header,
		Flags:              Flag(flags),
		FullCollectionName: collectionName,
		NumberToSkip:       numberToSkip,
		NumberToReturn:     numberToReturn,
		Query:              query,
	}

	return op, nil
}

// GetCollectionName returns the query collection name.
func (op *Query) GetCollectionName() string {
	return op.FullCollectionName
}

// IsCollection returns true when the specified name equals the full collection name, otherwise false.
func (op *Query) IsCollection(name string) bool {
	return name == op.FullCollectionName
}

// GetQuery returns the query document.
func (op *Query) GetQuery() bson.Document {
	return op.Query
}

// Documents returns the BSON documents.
func (op *Query) Documents() []bson.Document {
	return []bson.Document{op.Query}
}

// Size returns the message size including the header.
func (op *Query) Size() int32 {
	bodySize := 4 + (len(op.FullCollectionName) + 1) + 4 + 4 + len(op.Query)
	return int32(HeaderSize + bodySize)
}

// String returns the string description.
func (op *Query) String() string {
	return fmt.Sprintf("%s %X %s %d %d %s",
		op.Header.String(),
		op.Flags,
		op.FullCollectionName,
		op.NumberToSkip,
		op.NumberToReturn,
		op.Query.String(),
	)
}
