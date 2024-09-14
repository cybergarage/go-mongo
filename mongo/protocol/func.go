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
//
// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package protocol

import (
	"github.com/cybergarage/go-mongo/mongo/bson"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// AppendByte appends the byte to the buffer.
func AppendByte(dst []byte, val byte) []byte {
	return append(dst, val)
}

// AppendInt32 appends the int32 value to the buffer.
func AppendInt32(dst []byte, val int32) []byte {
	return bsoncore.AppendInt32(dst, val)
}

// AppendInt64 appends the int64 value to the buffer.
func AppendInt64(dst []byte, val int64) []byte {
	return bsoncore.AppendInt64(dst, val)
}

// AppendCString appends the string value to the buffer.
func AppendCString(dst []byte, val string) []byte {
	b := dst
	b = append(b, val...)
	return append(b, 0x00)
}

// AppendDocument appends the document to the buffer.
func AppendDocument(dst []byte, doc bson.Document) []byte {
	return append(dst, doc...)
}

// ReadInt32 reads a int32 value from src.
func ReadInt32(src []byte) (int32, []byte, bool) {
	if len(src) < 4 {
		return 0, src, false
	}
	return (int32(src[0]) | int32(src[1])<<8 | int32(src[2])<<16 | int32(src[3])<<24), src[4:], true
}

// ReadUint32 reads an uint32 value from src.
func ReadUint32(src []byte) (uint32, []byte, bool) {
	if len(src) < 4 {
		return 0, src, false
	}
	return (uint32(src[0]) | uint32(src[1])<<8 | uint32(src[2])<<16 | uint32(src[3])<<24), src[4:], true
}

// ReadInt64 reads a int64 value from src.
func ReadInt64(src []byte) (int64, []byte, bool) {
	return bsoncore.ReadInt64(src)
}

// ReadCString reads the cstring from src.
func ReadCString(src []byte) (string, []byte, bool) {
	return bsoncore.ReadKey(src)
}

// ReadCursorIDs reads numIDs cursor IDs from src.
func ReadCursorIDs(src []byte, numIDs int32) ([]int64, []byte, bool) {
	cursorIDs := []int64{}
	for i := 0; i < int(numIDs); i++ {
		id, src, ok := ReadInt64(src)
		if !ok {
			return cursorIDs, src, false
		}
		cursorIDs = append(cursorIDs, id)
	}
	return cursorIDs, src, true
}

// ReadSectionType reads the section type from src.
func ReadSectionType(src []byte) (SectionType, []byte, bool) {
	if len(src) < 1 {
		return SectionType(0), src, false
	}
	return SectionType(src[0]), src[1:], true
}

// ReadDocument reads a single document from src.
func ReadDocument(src []byte) (bson.Document, []byte, bool) {
	return bsoncore.ReadDocument(src)
}

// ReadDocuments reads as many documents as possible from src.
func ReadDocuments(src []byte) ([]bson.Document, []byte, bool) {
	docs := []bson.Document{}
	for {
		var doc bsoncore.Document
		var ok bool
		doc, src, ok = bsoncore.ReadDocument(src)
		if !ok {
			break
		}
		docs = append(docs, doc)
	}
	return docs, src, true
}

// ReadDocumentSequence reads an identifier and document sequence from src.
func ReadDocumentSequence(src []byte) (string, []bsoncore.Document, []byte, bool) {
	length, rem, ok := ReadInt32(src)
	if !ok || int(length) > len(src) {
		return "", nil, rem, false
	}

	rem, ret := rem[:length-4], rem[length-4:] // reslice so we can just iterate a loop later

	var identifier string
	identifier, rem, ok = ReadCString(rem)
	if !ok {
		return "", nil, rem, false
	}

	docs := make([]bsoncore.Document, 0)
	var doc bsoncore.Document
	for {
		doc, rem, ok = bsoncore.ReadDocument(rem)
		if !ok {
			break
		}
		docs = append(docs, doc)
	}
	if len(rem) > 0 {
		return "", nil, append(rem, ret...), false
	}

	return identifier, docs, ret, true
}
