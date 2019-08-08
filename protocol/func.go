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
	"github.com/cybergarage/go-mongo/bson"
	"go.mongodb.org/mongo-driver/x/mongo/driver/wiremessage"
)

// AppendByte appends the byte to the buffer.
func AppendByte(dst []byte, val byte) []byte {
	return wiremessage.AppendMsgSectionType(dst, wiremessage.SectionType(val))
}

// AppendInt32 appends the int32 value to the buffer.
func AppendInt32(dst []byte, val int32) []byte {
	return wiremessage.AppendKillCursorsNumberIDs(dst, val)
}

// AppendInt64 appends the int64 value to the buffer.
func AppendInt64(dst []byte, val int64) []byte {
	return wiremessage.AppendReplyCursorID(dst, val)
}

// AppendCString appends the string value to the buffer.
func AppendCString(dst []byte, val string) []byte {
	return wiremessage.AppendQueryFullCollectionName(dst, val)
}

// AppendDocument appends the document to the buffer.
func AppendDocument(dst []byte, doc bson.Document) []byte {
	return append(dst, doc...)
}

// ReadInt32 reads a int32 value from src.
func ReadInt32(src []byte) (int32, []byte, bool) {
	return wiremessage.ReadReplyStartingFrom(src)
}

// ReadUint32 reads an uint32 value from src.
func ReadUint32(src []byte) (uint32, []byte, bool) {
	val, rem, ok := ReadInt32(src)
	return uint32(val), rem, ok
}

// ReadInt64 reads a int64 value from src.
func ReadInt64(src []byte) (int64, []byte, bool) {
	return wiremessage.ReadReplyCursorID(src)
}

// ReadCString reads the cstring from src.
func ReadCString(src []byte) (string, []byte, bool) {
	return wiremessage.ReadQueryFullCollectionName(src)
}

// ReadDocument reads a single document from src.
func ReadDocument(src []byte) (bson.Document, []byte, bool) {
	return wiremessage.ReadReplyDocument(src)
}

// ReadDocuments reads as many documents as possible from src
func ReadDocuments(src []byte) ([]bson.Document, []byte, bool) {
	return wiremessage.ReadReplyDocuments(src)
}
