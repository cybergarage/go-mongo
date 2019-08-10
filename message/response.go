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

package message

import (
	"strconv"

	"github.com/cybergarage/go-mongo/bson"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

const (
	ok         = "ok"
	cursor     = "cursor"
	firstBatch = "firstBatch"
	nameSpace  = "ns"
)

// Response represents response elements
type Response struct {
	*bson.Dictionary
}

// NewResponse returns a new response instance.
func NewResponse() *Response {
	res := &Response{
		Dictionary: bson.NewDictionary(),
	}
	return res
}

// NewResponseWithElements returns a new response instance.
func NewResponseWithElements(elements map[string]interface{}) *Response {
	res := NewResponse()
	res.SetElements(elements)
	return res
}

// SetStatus sets an int32 response result.
func (res *Response) SetStatus(flag bool) {
	if flag {
		res.SetDoubleElement(ok, 1.0)
		return
	}
	res.SetDoubleElement(ok, 0.0)
}

// SetCursorDocuments sets a resultset.
func (res *Response) SetCursorDocuments(fullCollectionName string, docs []bson.Document) {
	var arrIdx int32
	cursorIdx, cursorDoc := bsoncore.AppendDocumentStart(nil)
	arrIdx, cursorDoc = bsoncore.AppendArrayElementStart(cursorDoc, firstBatch)
	for n, doc := range docs {
		cursorDoc = bsoncore.AppendDocumentElement(cursorDoc, strconv.Itoa(n), doc)
	}
	cursorDoc, _ = bsoncore.AppendArrayEnd(cursorDoc, arrIdx)

	cursorDoc = bsoncore.AppendInt64Element(cursorDoc, "id", 0)
	cursorDoc = bsoncore.AppendStringElement(cursorDoc, nameSpace, fullCollectionName)
	cursorDoc, _ = bsoncore.AppendDocumentEnd(cursorDoc, cursorIdx)

	res.SetDocumentElement(cursor, cursorDoc)
}
