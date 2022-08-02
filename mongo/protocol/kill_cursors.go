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

// KillCursors represents a OP_KILL_CURSORS of MongoDB wire protocol.
// See : MongoDB Wire Protocol
// https://docs.mongodb.com/manual/reference/mongodb-wire-protocol/
type KillCursors struct {
	*Header                   // A standard wire protocol header
	ZERO              int32   // 0 - reserved for future use
	NumberOfCursorIDs int32   // number of documents to return
	CursorIDs         []int64 // cursorID from the OP_REPLY
}

// NewKillCursorsWithHeaderAndBody returns a new get more instance with the specified bytes.
func NewKillCursorsWithHeaderAndBody(header *Header, body []byte) (*KillCursors, error) {
	zero, offsetBody, ok := ReadInt32(body)
	if !ok {
		return nil, newMessageRequestError(OpDelete, body)
	}

	numberOfCursorIDs, offsetBody, ok := ReadInt32(offsetBody)
	if !ok {
		return nil, newMessageRequestError(OpKillCursors, body)
	}

	cursorIDs, _, ok := ReadCursorIDs(offsetBody, numberOfCursorIDs)
	if !ok {
		return nil, newMessageRequestError(OpKillCursors, body)
	}

	op := &KillCursors{
		Header:            header,
		ZERO:              zero,
		NumberOfCursorIDs: numberOfCursorIDs,
		CursorIDs:         cursorIDs,
	}

	return op, nil
}

// Size returns the message size including the header.
func (op *KillCursors) Size() int32 {
	bodySize := 4 + 4 + (8 * op.NumberOfCursorIDs)
	return HeaderSize + bodySize
}

// String returns the string description.
func (op *KillCursors) String() string {
	str := fmt.Sprintf("%s %d ",
		op.Header.String(),
		op.NumberOfCursorIDs,
	)

	for _, cursorID := range op.CursorIDs {
		str += fmt.Sprintf("%X ", cursorID)
	}

	return str
}
