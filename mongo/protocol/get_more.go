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

// GetMore represents a OP_GET_MORE of MongoDB wire protocol.
// See : MongoDB Wire Protocol
// https://docs.mongodb.com/manual/reference/mongodb-wire-protocol/
type GetMore struct {
	*Header                   // A standard wire protocol header
	ZERO               int32  // 0 - reserved for future use
	FullCollectionName string // "dbname.collectionname"
	NumberToReturn     int32  // number of documents to return
	CursorID           int64  // cursorID from the OP_REPLY
}

// NewGetMoreWithHeaderAndBody returns a new get more instance with the specified bytes.
func NewGetMoreWithHeaderAndBody(header *Header, body []byte) (*GetMore, error) {
	zero, offsetBody, ok := ReadInt32(body)
	if !ok {
		return nil, newMessageRequestError(OpDelete, body)
	}

	collectionName, offsetBody, ok := ReadCString(offsetBody)
	if !ok {
		return nil, newMessageRequestError(OpGetMore, body)
	}

	numberToReturn, offsetBody, ok := ReadInt32(offsetBody)
	if !ok {
		return nil, newMessageRequestError(OpGetMore, body)
	}

	cursorID, _, ok := ReadInt64(offsetBody)
	if !ok {
		return nil, newMessageRequestError(OpGetMore, body)
	}

	op := &GetMore{
		Header:             header,
		ZERO:               zero,
		FullCollectionName: collectionName,
		NumberToReturn:     numberToReturn,
		CursorID:           cursorID,
	}

	return op, nil
}

// Size returns the message size including the header.
func (op *GetMore) Size() int32 {
	bodySize := 4 + (len(op.FullCollectionName) + 1) + 4 + 8
	return int32(HeaderSize + bodySize)
}

// String returns the string description.
func (op *GetMore) String() string {
	return fmt.Sprintf("%s %s %d %d",
		op.Header.String(),
		op.FullCollectionName,
		op.NumberToReturn,
		op.CursorID,
	)
}
